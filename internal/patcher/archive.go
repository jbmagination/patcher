package patcher

import (
	"compress/flate"
	"log"
	"os"

	"github.com/mholt/archiver"
)

// Extract Discord's IPA
func extractDiscord(discordPath *string) {
	log.Println("extracting", *discordPath)
	format := archiver.Zip{}

	err := format.Unarchive(*discordPath, ".")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(*discordPath, "extracted")
}

// Pack Discord's IPA
func packDiscord() {
	log.Println("packing discord")

	format := archiver.Zip{
		CompressionLevel: flate.BestCompression,
	}
	err := format.Archive([]string{"Payload"}, "Discord.zip")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Rename("Discord.zip", "Enmity.ipa")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Discord packed")
}
