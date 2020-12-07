package config

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)


func TestNewConfig(t *testing.T) {

	badfile, err := ioutil.TempFile("", "conf.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(badfile.Name())
	badfile.WriteString(`aefSD
sadfg
RFABND FYGUMG
V`)
	badfile.Sync()

	goodfile, err := ioutil.TempFile("", "conf.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(goodfile.Name())
	goodfile.WriteString(`[storage]
inMemory = true
SQLHost = "localhost"`)
	goodfile.Sync()

	t.Run("No such file", func(t *testing.T) {
		var c Calendar
		e := New("adfergdth", &c)
		require.Equal(t, Calendar{}, c)
		require.Error(t, e)
	})

	t.Run("Bad file", func(t *testing.T) {
		var c Calendar
		e := New(badfile.Name(), &c)
		require.Equal(t, Calendar{}, c)
		require.Error(t, e)
	})

	t.Run("TOML reading", func(t *testing.T) {
		var c Calendar
		e := New(goodfile.Name(), &c)
		require.Equal(t, true, c.Storage.InMemory)
		require.Equal(t, "localhost", c.Storage.SQLHost)
		require.NoError(t, e)
	})

	t.Run("ENV reading", func(t *testing.T) {
		for k, v := range map[string]string{"APP_STRUCT1_VAR1": "val1", "APP_STRUCT1_VAR2": "val2", "APP_STRUCT2_VAR1": "val3", "APP_STRUCT2_VAR2": "val4", "APP_STRUCT3_VAR1": "val5", "APP_STRUCT3_VAR2": "val6"} {
			require.NoError(t, os.Setenv(k, v))
		}
		var str struct {
			Struct1 struct {
				Var1 string
				Var2 string
			}
			Struct2 struct {
				Var1 string
				Var2 string
			}
			Struct3 struct {
				Var1 string
				Var2 string
			}
		}

		err := New("", &str)
		require.NoError(t, err)
		require.Equal(t, "val1", str.Struct1.Var1)
		require.Equal(t, "val2", str.Struct1.Var2)
		require.Equal(t, "val3", str.Struct2.Var1)
		require.Equal(t, "val4", str.Struct2.Var2)
		require.Equal(t, "val5", str.Struct3.Var1)
		require.Equal(t, "val6", str.Struct3.Var2)
	})

}
