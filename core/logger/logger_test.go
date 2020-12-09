package logger

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	oslog "log"
	"os"
	"strings"
	"testing"
)

func TestLoggerLogic(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "log.")
	if err != nil {
		oslog.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	conf := Config{File: tmpfile.Name(), Level: "warn", MuteStdout: false}
	log, err := New(conf)
	if err != nil {
		oslog.Fatal(err)
	}

	t.Run("Messages arround the level", func(t *testing.T) {
		log.Debugf("debug message")
		log.Errorf("error message")

		res, err := ioutil.ReadAll(tmpfile)
		if err != nil {
			oslog.Fatal(err)
		}
		require.Less(t, strings.Index(string(res), "debug message"), 0)
		require.Greater(t, strings.Index(string(res), "error message"), 0)
	})
}

func TestLoggerNegative(t *testing.T) {
	t.Run("Bad file name", func(t *testing.T) {
		conf := Config{File: "", Level: "debug", MuteStdout: true}
		_, err := New(conf)
		require.Error(t, err, "invalid logger config")
	})

	t.Run("Bad level", func(t *testing.T) {
		conf := Config{File: "asdafad", Level: "wegretryjt", MuteStdout: true}
		_, err := New(conf)
		require.Error(t, err, "invalid logger config")
	})
}
