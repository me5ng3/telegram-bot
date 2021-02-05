package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jessevdk/go-flags"
)

var (
	regexMAC = regexp.MustCompile(`^([0-9a-fA-F]{2}[` + ":-" + `]){5}([0-9a-fA-F]{2})$`)

	cliFlags struct {
		Help        bool   `short:"h" long:"help"`
		Interface   string `short:"i" long:"interface" default:""`
		BroadcastIP string `short:"b" long:"bcast" default:"255.255.255.255"`
		UDPPort     string `short:"p" long:"port" default:"9"`
	}
)

type MACAdress [6]byte

type MagicPacket struct {
	header  [6]byte
	payload [16]MACAdress
}

func newMagicPacket(address string) (*MagicPacket, error) {
	var packet MagicPacket
	var add MACAdress

	hwAddr, err := net.ParseMAC(address)
	if err != nil {
		return nil, err
	}

	if !regexMAC.MatchString(address) {
		return nil, fmt.Errorf("Wrong MAC address format.")
	}

	for b := range add {
		add[b] = hwAddr[b]
	}

	for b := range packet.header {
		packet.header[b] = 0xFF
	}

	for b := range packet.payload {
		packet.payload[b] = add
	}

	return &packet, nil
}

func (packet *MagicPacket) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, packet); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func sendWakeSignal(macAdress string, args []string) error {
	if cliFlags.BroadcastInterface != "" {
		bcastInterface = cliFlags.BroadcastInterface
	}

	var localAddr *net.UDPAddr
	if bcastInterface != "" {
		localAddr, err = ipFromInterface(bcastInterface)
		if err != nil {
			return err
		}
	}

	bcastAddr := fmt.Sprintf("%s:%s", cliFlags.BroadcastIP, cliFlags.UDPPort)
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	magicPacket, err := newMagicPacket(macAdress)
	if err != nil {
		return err
	}

	b, err := magicPacket.Marshal()
	if err != nil {
		return err
	}

	req, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer req.Close()

	n, err := req.Write(b)
	if err == nil && n != 102 {
		return fmt.Errorf("Magic Packet didn't send all the bytes. Only sent %d bytes.", n)
	}
	if err != nil {
		return err
	}
}

func (cmdHandler *commandHandler) wol(u *telegram.Update) {
	msg := strings.Split(telegram.NewMessage(u.Message.Chat.ID, u.Message.Text).Text, " ")

	parser := flags.NewParser(&cliFlags, flags.Default & ^flags.HelpFlag)
	args, err := parser.Parse()

	switch {
	case err != nil:
		cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "Could not parse message. Check logs for more info."))
		cmdHandler.logger.consumer <- err
		break

	case len(args) == 0:
		cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "Please specify the MAC address and the port to send the Wake On Lan Magic Packet to. Check --help for more informations."))
		break

	case true:
		cmd, cmdArgs := strings.ToLower(args[0]), args[1:]
		if fn, ok := cmdMap[cmd]; ok {
			err = fn(cmdArgs, aliases)
		} else {
			err = sendWakeSignal(args, aliases)
			if err != nil {
				cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "Could not send wake on lan packet."))
				cmdHandler.logger.consumer <- err
				break
			}
			cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "Wake On Lan packet sent successfully."))
			cmdHandler.log.consumer <- "Wake On Lan Packet sent successfully"
		}
	}
}
