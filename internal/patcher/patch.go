package patcher

import (
	"errors"
	"log"
	"os"

	"howett.net/plist"
)

const (
	DEFAULT_IPA_PATH    = "files/Discord.ipa"
)

const (
	IPA_URL    = "https://github.com/enmity-mod/tweak/blob/main/Discord.ipa?raw=true"
)

func PatchDiscord(discordPath *string, iconsPath *string, dylibPath *string) {
	log.Println("starting patcher")

	checkFile(discordPath, DEFAULT_IPA_PATH, IPA_URL)
	checkFile(iconsPath, DEFAULT_ICONS_PATH, ICONS_URL)

	extractDiscord(discordPath)

	log.Println("adding Enmity url scheme")
	if err := patchSchemes(); err != nil {
		log.Fatalln(err)
	}
	log.Println("url scheme added")

	log.Println("remove devices whitelist")
	if err := patchDevices(); err != nil {
		log.Fatalln(err)
	}
	log.Println("device whitelist removed")

	log.Println("showing Discord's document folder in the Files app and Finder/iTunes")
	if err := patchiTunesAndFiles(); err != nil {
		log.Fatalln(err)
	}
	log.Println("patched")

	packDiscord()
	log.Println("cleaning up")
	clearPayload()

	log.Println("done!")
}

// Check if file exists, download if not found
func checkFile(path *string, defaultPath string, url string) {
	_, err := os.Stat(*path)
	if errors.Is(err, os.ErrNotExist) {
		if *path == defaultPath {
			log.Println("downloading", url, "to", *path)
			err := downloadFile(url, path)
			if err != nil {
				log.Println("error downloading", url)
				log.Fatalln(err)
			}
		} else {
			log.Fatalln("file not found", *path)
		}
	}
}

// Delete the payload folder
func clearPayload() {
	err := os.RemoveAll("Payload")
	if err != nil {
		log.Panicln(err)
	}
}

// Load Discord's plist file
func loadPlist() (map[string]interface{}, error) {
	infoFile, err := os.Open("Payload/Discord.app/Info.plist")
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	decoder := plist.NewDecoder(infoFile)
	if err := decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}

// Save Discord's plist file
func savePlist(info *map[string]interface{}) error {
	infoFile, err := os.OpenFile("Payload/Discord.app/Info.plist", os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	encoder := plist.NewEncoder(infoFile)
	err = encoder.Encode(*info)
	return err
}

// Patch Discord's URL scheme to add Enmity's URL handler
func patchSchemes() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	info["CFBundleURLTypes"] = append(
		info["CFBundleURLTypes"].([]interface{}),
		map[string]interface{}{
			"CFBundleURLName": "Enmity",
			"CFBundleURLSchemes": []string{
				"enmity",
			},
		},
	)

	err = savePlist(&info)
	return err
}

// Remove Discord's device limits
func patchDevices() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	delete(info, "UISupportedDevices")

	err = savePlist(&info)
	return err
}

// Show Enmity's document folder in Files app and iTunes/Finder
func patchiTunesAndFiles() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}
	info["UISupportsDocumentBrowser"] = true
	info["UIFileSharingEnabled"] = true

	err = savePlist(&info)
	return err
}
