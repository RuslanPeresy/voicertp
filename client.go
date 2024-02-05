package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
)

func createPionRTPPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version: 2,
			// Taken from Discord voice docs
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}
func handleVoice(c chan *discordgo.Packet, conn net.Conn) {
	buffer := make([]byte, 1500)
	for p := range c {
		// Construct pion RTP packet from DiscordGo's type.
		rtpPacket := createPionRTPPacket(p)

		fmt.Println(rtpPacket)
		// Marshal into buffer
		n, err := rtpPacket.MarshalTo(buffer)
		if err != nil {
			panic(err)
		}

		// Write
		if _, writeErr := conn.Write(buffer[:n]); writeErr != nil {
			// For this particular example, third party applications usually timeout after a short
			// amount of time during which the user doesn't have enough time to provide the answer
			// to the browser.
			// That's why, for this particular example, the user first needs to provide the answer
			// to the browser then open the third party application. Therefore we must not kill
			// the forward on "connection refused" errors
			var opError *net.OpError
			if errors.As(writeErr, &opError) && opError.Err.Error() == "write: connection refused" {
				continue
			}
			panic(err)
		}
	}
}

func main() {
	var Token = ""     // Set bot token
	var ChannelID = "" // Set channel id
	var GuildID = ""   // Set server id

	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
		return
	}
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", 8080))
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	defer func(conn net.PacketConn) {
		if closeErr := conn.Close(); closeErr != nil {
			panic(closeErr)
		}
	}(conn)

	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	defer s.Close()

	// We only really care about receiving voice state updates.
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	v, err := s.ChannelVoiceJoin(GuildID, ChannelID, true, false)
	if err != nil {
		fmt.Println("failed to join voice channel:", err)
		return
	}

	handleVoice(v.OpusRecv, conn)
}
