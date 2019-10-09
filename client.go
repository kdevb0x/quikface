// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	_ "github.com/gobwas/ws"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/pion/webrtc/pkg/media"
	rtc "github.com/pion/webrtc/v2"
)

type Client struct {
	// uuid
	ID          uint32
	DisplayName string
	// Client network address
	Addr string
	// for signal mesages; known as "rwc" in collider
	// TODO: Maybe change this to io.ReadWriteCloser and have signalFunc
	// implement the interface?
	Signal io.ReadWriteCloser // signalFunc
	// Signal io.ReadWriteCloser

	MsgQueue []Message
}

// randomChatName returns a random displayname using simple {adj}{noun} lists
// (ala demonsaw; thx for idea chudz!)
func randomChatName() string {
	var list1 = []string{"Adaptable",
		"Adventurous",
		"Affable",
		"Affectionate",
		"Agreeable",
		"Ambitious",
		"Amiable",
		"Amicable",
		"Amusing",
		"Brave",
		"Bright",
		"Calm",
		"Careful",
		"Charming",
		"Communicative",
		"Compassionate",
		"Conscientious",
		"Considerate",
		"Convivial",
		"Courageous",
		"Courteous",
		"Creative",
		"Decisive",
		"Determined",
		"Diligent",
		"Diplomatic",
		"Discreet",
		"Dynamic",
		"Easygoing",
		"Emotional",
		"Energetic",
		"Enthusiastic",
		"Exuberant",
		"Faithful",
		"Fearless",
		"Forceful",
		"Frank",
		"Friendly",
		"Funny",
		"Generous",
		"Gentle",
		"Good",
		"Gregarious",
		"Helpful",
		"Honest",
		"Humorous",
		"Imaginative",
		"Impartial",
		"Independent",
		"Intellectual",
		"Intelligent",
		"Intuitive",
		"Inventive",
		"Kind",
		"Loving",
		"Loyal",
		"Modest",
		"Neat",
		"Nice",
		"Optimistic",
		"Passionate",
		"Patient",
		"Persistent",
		"Pioneering",
		"Philosophical",
		"Placid",
		"Plucky",
		"Polite",
		"Powerful",
		"Practical",
		"Quiet",
		"Rational",
		"Reliable",
		"Reserved",
		"Resourceful",
		"Romantic",
		"Sensible",
		"Sensitive",
		"Shy",
		"Sincere",
		"Sociable",
		"Straightforward",
		"Sympathetic",
		"Thoughtful",
		"Tidy",
		"Tough",
		"Unassuming",
		"Understanding",
		"Versatile",
		"Warmhearted",
		"Willing",
		"Witty"}

	var list2 = []string{"Aardvark",
		"Aardwolf",
		"Elephant",
		"Pangolin",
		"Alligator",
		"Alpaca",
		"Anteater",
		"Antelope",
		"Ape(s)",
		"Horse",
		"Armadillo",
		"Baboon",
		"Badger",
		"Bandicoot",
		"Tiger",
		"Beaver",
		"Whale",
		"Goat",
		"Bison",
		"Rhino",
		"Monkey",
		"Boar",
		"Bobcat",
		"Bonobo",
		"Dolphin",
		"Buffalo",
		"Bull",
		"Camel",
		"Capybara",
		"Caribou",
		"Cattle",
		"Cheetah",
		"Chimpanzee",
		"Chinchilla",
		"Chipmunk",
		"Seal",
		"Cougar",
		"Coyote",
		"Crocodile",
		"Frog",
		"Deer",
		"Degus",
		"Dingo",
		"Dolphin",
		"Donkey",
		"Dormouse",
		"Dugong",
		"Elephant",
		"Elk",
		"Ermine",
		"Lynx",
		"Ferret",
		"Panther",
		"Fox",
		"Frog",
		"Tortoise",
		"Gazelle",
		"Gecko",
		"Gibbon",
		"Giraffe",
		"Goat",
		"Gopher",
		"Gorilla",
		"Groundhog",
		"Hare",
		"Hedgehog",
		"Hippopotamus",
		"Hyena",
		"Hyrax",
		"Iguana",
		"Iguanodon",
		"Impala",
		"Jackal",
		"Jackrabbit",
		"Jaguar",
		"Jellyfish",
		"Kangaroo",
		"Koala",
		"Kookaburra",
		"Lama",
		"Lamb",
		"Lancelet",
		"Lemming",
		"Lemur",
		"Leopard",
		"Lion",
		"Llama",
		"Lynx",
		"Manatee",
		"Mantis",
		"Marmot",
		"Meerkat",
		"Mink",
		"Mole",
		"Mongoose",
		"Monkey",
		"Moose",
		"Mouse",
		"Mule",
		"Muskox",
		"Muskrat",
		"Narwhal",
		"Nautilus",
		"Newt",
		"Nutria",
		"Nyala",
		"Ocelot",
		"Octopus",
		"Okapi",
		"Opossum",
		"Orangutan",
		"Orca",
		"Otter",
		"Ox",
		"Panda",
		"Panther",
		"Pig",
		"Platypus",
		"Porcupine",
		"Porpoise",
		"Possum",
		"Potto",
		"Puma",
		"Quokkas",
		"Quolls",
		"Rabbit",
		"Raccoon",
		"Rat",
		"Ray",
		"Reindeer",
		"Rhino",
		"Rhinoceros",
		"Salamander",
		"Seal",
		"Shark",
		"Sheep",
		"Skink",
		"Skunk",
		"Sloth",
		"Squirrel",
		"Takin",
		"Tamarin",
		"Tapir",
		"Terrapin",
		"Tiger",
		"Topi",
		"Tortoise",
		"Turtle",
		"Uakari",
		"Vicuna",
		"Vole",
		"Wallaby",
		"Walrus",
		"Warthog",
		"Weasel",
		"Wildcat",
		"Wildebeest",
		"Wolf",
		"Wolverine",
		"Wombat",
		"Woodchuck",
		"Yak",
		"Zebra",
		"Zebu",
		"Zorilla"}

	rand.Seed(time.Now().UnixNano())
	l1pick := rand.Intn(len(list1))
	l2pick := rand.Intn(len(list2))
	return list1[l1pick] + " " + list2[l2pick]
}

func NewClient(displayname ...string) *Client {
	if len(displayname) > 0 {
		// cant break if stmnts, so had to flip this test
		if displayname[0] != "" {
			return &Client{ID: uuid.New().ID(), DisplayName: displayname[0], MsgQueue: make([]Message, 0, 10)}

		}
		// here displayname[0] == ""
	}
	return &Client{ID: uuid.New().ID(), DisplayName: randomChatName(), MsgQueue: make([]Message, 0, 10)}

}

func (c *Client) Register(signaler io.ReadWriteCloser) error {
	if c.Signal != nil {
		return fmt.Errorf("duplicate registration; %s already has a signal connection registered", c.ID)
	}
	c.Signal = signaler
	return nil
}

// JoinRoom return a pointer to the room if the join was successfull, otherwise
// said pointer is nil, and err explains why.
func (c *Client) JoinRoom(name string, masterDirectory *RoomList) (*Room, error) {
	if hash, exists := masterDirectory.Rooms[name]; exists {
		room, err := masterDirectory.GetRoom(hash)
		if err != nil {
			return nil, fmt.Errorf("error finding hash of %s, %w\n", name, err)
		}
		room.Clients[c.ID] = c
		return room, nil

	}
	return nil, fmt.Errorf("error: room %s doesn't exist", name)
}

func (c *Client) initWebRTCSession(signal *websocket.Conn) (*rtc.PeerConnection, error) {

}
