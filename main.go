
package main

import (
	"context"
	"encoding/json"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	token := "" // use os.Getenv("TOKEN")
	vk := api.NewVK(token)

	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)



		if len(obj.Message.Attachments)  != 0 {
			log.Println(obj.Message.Attachments[0].Photo.Sizes[4])
			var phUrl string = obj.Message.Attachments[0].Photo.Sizes[4].URL

			response1, e := http.Get(phUrl)
			if e != nil {
				log.Fatal(e)
			}
			defer response1.Body.Close()
			var fileName = ("C:\\Users\\YaTeb\\go\\src\\VkBot\\tmp\\" + strconv.Itoa(obj.Message.PeerID) + ".jpg")
			file, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(file, response1.Body)
			if err != nil {
				log.Fatal(err)
			}

			photo1, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer photo1.Close()

			type Params map[string]interface{}
			uploadServer, err := vk.PhotosGetMessagesUploadServer(api.Params(Params{
				"peer_id": obj.Message.PeerID,
			}))
			if err != nil {
				return
			}
			bodyContent, err := vk.UploadFile(uploadServer.UploadURL,  photo1, "photo", "photo.jpeg")
			if err != nil {
				return
			}

			var handler object.PhotosMessageUploadResponse


			err = json.Unmarshal(bodyContent, &handler)
			if err != nil {
				return
			}
			response, err := vk.PhotosSaveMessagesPhoto(api.Params(Params{
				"server": handler.Server,
				"photo":  handler.Photo,
				"hash":   handler.Hash,
			}))


			b := params.NewMessagesSendBuilder()
			b.Message("Вот твоя картиночка")
			b.RandomID(2)
			b.Attachment(response)
			b.PeerID(obj.Message.PeerID)


			_, errSendPh := vk.MessagesSend(b.Params)
			if errSendPh != nil {
				log.Fatal(err)
			}
		}

		if obj.Message.FromID ==206312673 {
			b := params.NewMessagesSendBuilder()
			b.Message("Денис пидор")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)


			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}
	})




	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}