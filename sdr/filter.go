package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// sensorFilter is the allow-list of sensors, the stand-in for the console's
// ID learning: the console registers the first ID it hears per sensor type
// and ignores every other one, so a neighbour's station cannot leak in.
//
// Entries are "<model>/<id>", the same pair that forms the MQTT topic. An
// empty list accepts everything.
type sensorFilter struct {
	allowed map[string]bool
	seen    map[string]bool // sensors already mentioned in the log
}

func sensorKey(model string, id int) string { return model + "/" + strconv.Itoa(id) }

// parseSensors reads a comma- or space-separated list of "<model>/<id>".
func parseSensors(list string) (*sensorFilter, error) {
	f := &sensorFilter{allowed: map[string]bool{}, seen: map[string]bool{}}
	for _, e := range strings.FieldsFunc(list, func(r rune) bool { return r == ',' || r == ' ' || r == '\n' || r == '\t' }) {
		model, id, ok := strings.Cut(e, "/")
		n, err := strconv.Atoi(id)
		if !ok || model == "" || err != nil {
			return nil, fmt.Errorf("bad sensor %q: want <model>/<id>, e.g. Fineoffset-WS80/34306", e)
		}
		f.allowed[sensorKey(model, n)] = true
	}
	return f, nil
}

// Allow reports whether a message should be output. The first time a sensor
// is heard it is logged, so the list can be built from the log.
func (f *sensorFilter) Allow(m *Message) bool {
	key := sensorKey(m.Model, m.ID)
	ok := len(f.allowed) == 0 || f.allowed[key]
	if !f.seen[key] {
		f.seen[key] = true
		if ok {
			log.Printf("sensor %s: receiving", key)
		} else {
			log.Printf("sensor %s: not in SDR_SENSORS, ignoring it", key)
		}
	}
	return ok
}
