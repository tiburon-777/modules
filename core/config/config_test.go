package config

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type TestConf struct {
	Section1 struct {
		VarInt1    int
		VarString1 string
		VarBool1   bool
	}
	Section2 struct {
		VarInt2    int
		VarString2 string
		VarBool2   bool
	}
}

func TestSetFromFilePositive(t *testing.T) {

	goodfile, err := ioutil.TempFile("", "conf.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(goodfile.Name())
	goodfile.WriteString(
		`[section1]
				varint1 = 11
				varstring1 = "first string"
				varbool1 = true
			[section2]
				varint2 = 22
				varstring2 = "second string"
				varbool2 = true`)
	goodfile.Sync()

	partfile, err := ioutil.TempFile("", "conf.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(partfile.Name())
	partfile.WriteString(
		`[section1]
				varint1 = 11
			[section2]
				varbool2 = true`)
	partfile.Sync()

	zerofile, err := ioutil.TempFile("", "conf.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(zerofile.Name())
	zerofile.WriteString(
		`[sdgsergergse]
				argargragffg = "sdgfhdtn"
			[sdbgjuykb]
				fbfjmuinuyhg = 134134`)
	zerofile.Sync()

	// Если файл валидный TOML, метод вернет заполненный конфиг и nil.
	t.Run("Successful parsing the file", func(t *testing.T) {
		var c TestConf
		i := New(&c)
		err := i.SetFromFile(goodfile.Name())
		require.NoError(t, err)
		require.Equal(t, 11, c.Section1.VarInt1)
		require.Equal(t, "first string", c.Section1.VarString1)
		require.Equal(t, true, c.Section1.VarBool1)
		require.Equal(t, 22, c.Section2.VarInt2)
		require.Equal(t, "second string", c.Section2.VarString2)
		require.Equal(t, true, c.Section2.VarBool2)
	})

	// Если некоторых переменных нет в файле, метод вернет конфиг заполненный насколько возможно и nil.
	t.Run("Partial config applying from file", func(t *testing.T) {
		var c TestConf
		i := New(&c)
		err := i.SetFromFile(partfile.Name())
		require.NoError(t, err)
		require.Equal(t, 11, c.Section1.VarInt1)
		require.Equal(t, "", c.Section1.VarString1)
		require.Equal(t, false, c.Section1.VarBool1)
		require.Equal(t, 0, c.Section2.VarInt2)
		require.Equal(t, "", c.Section2.VarString2)
		require.Equal(t, true, c.Section2.VarBool2)
	})

	// Если переменных нет в файле, метод вернет zerovalue конфиг и nil.
	t.Run("No vars in file", func(t *testing.T) {
		var c TestConf
		i := New(&c)
		err := i.SetFromFile(zerofile.Name())
		require.NoError(t, err)
		require.Equal(t, TestConf{}, c)
	})
}

func TestSetFromFileNegative(t *testing.T) {

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

	typesfile, err := ioutil.TempFile("", "conf.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(typesfile.Name())
	typesfile.WriteString(
		`[section1]
				varint1 = "first string"
				varstring1 = true
				varbool1 = 11
			[section2]
				varint2 = 22
				varstring2 = "second string"
				varbool2 = true`)
	typesfile.Sync()

	// Если файла не существует, метод вернет пустой конфиг и ошибку.
	t.Run("File doesn't exist", func(t *testing.T) {
		var c TestConf
		i := New(&c)
		err := i.SetFromFile("dfsdfgsdfds")
		require.Error(t, err)
		require.Equal(t, TestConf{}, c)
	})

	// Если файл не является TOML-readable, метод вернет пустой конфиг и ошибку.
	t.Run("Try to use corrupted file", func(t *testing.T) {
		var c TestConf
		i := New(&c)
		err := i.SetFromFile(badfile.Name())
		require.Error(t, err)
		require.Equal(t, TestConf{}, c)
	})

	// Если некоторые типы перепутаны, метод вернет zerovalue конфиг и ошибку.
	t.Run("Unexpected types", func(t *testing.T) {
		var c TestConf
		i := New(&c)
		err := i.SetFromFile(typesfile.Name())
		require.Error(t, err)
		require.Equal(t, TestConf{}, c)
	})
}

func TestSetFromEnvPositive(t *testing.T) {

	// Если все соответствующие переменные есть в окружении, метод вернет полностью заполненный конфиг и nil.
	t.Run("Fulfill config applying from env", func(t *testing.T) {
		for k, v := range map[string]string{"APP_SECTION1_VARINT1": "11", "APP_SECTION1_VARSTRING1": "first string", "APP_SECTION1_VARBOOL1": "true", "APP_SECTION2_VARINT2": "22", "APP_SECTION2_VARSTRING2": "second string", "APP_SECTION2_VARBOOL2": "true"} {
			require.NoError(t, os.Setenv(k, v))
		}
		var c TestConf
		i := New(&c)
		err := i.SetFromEnv("APP")
		require.NoError(t, err)
		require.Equal(t, 11, c.Section1.VarInt1)
		require.Equal(t, "first string", c.Section1.VarString1)
		require.Equal(t, true, c.Section1.VarBool1)
		require.Equal(t, 22, c.Section2.VarInt2)
		require.Equal(t, "second string", c.Section2.VarString2)
		require.Equal(t, true, c.Section2.VarBool2)
	})

	// Если некоторых переменных нет в окружении, метод вернет конфиг заполненный насколько возможно и nil.
	t.Run("Partial config applying from env", func(t *testing.T) {
		for k, v := range map[string]string{"APP_SECTION1_VARINT1": "11", "APP_SECTION1_VARSTRING1": "first string", "APP_SECTION1_VARBOOL1": "true"} {
			require.NoError(t, os.Setenv(k, v))
		}
		for k, _ := range map[string]string{"APP_SECTION2_VARINT2": "22", "APP_SECTION2_VARSTRING2": "second string", "APP_SECTION2_VARBOOL2": "true"} {
			require.NoError(t, os.Unsetenv(k))
		}
		var c TestConf
		i := New(&c)
		err := i.SetFromEnv("APP")
		require.NoError(t, err)
		require.Equal(t, 11, c.Section1.VarInt1)
		require.Equal(t, "first string", c.Section1.VarString1)
		require.Equal(t, true, c.Section1.VarBool1)
		require.Equal(t, 0, c.Section2.VarInt2)
		require.Equal(t, "", c.Section2.VarString2)
		require.Equal(t, false, c.Section2.VarBool2)
	})

	// Если переменных нет в окружении, метод вернет zerovalue конфиг и nil.
	t.Run("No env vars", func(t *testing.T) {
		for k, _ := range map[string]string{"APP_SECTION1_VARINT1": "11", "APP_SECTION1_VARSTRING1": "first string", "APP_SECTION1_VARBOOL1": "true", "APP_SECTION2_VARINT2": "22", "APP_SECTION2_VARSTRING2": "second string", "APP_SECTION2_VARBOOL2": "true"} {
			require.NoError(t, os.Unsetenv(k))
		}
		var c TestConf
		i := New(&c)
		err := i.SetFromEnv("APP")
		require.NoError(t, err)
		require.Equal(t, TestConf{}, c)
	})
}

func TestSetFromEnvNegative(t *testing.T) {

	// Если некоторые типы перепутаны, метод вернет zerovalue конфиг и ошибку.
	t.Run("Unexpected types", func(t *testing.T) {
		for k, v := range map[string]string{"APP_SECTION1_VARINT1": "first string", "APP_SECTION1_VARSTRING1": "false", "APP_SECTION1_VARBOOL1": "shstrjerthgccw", "APP_SECTION2_VARINT2": "22", "APP_SECTION2_VARSTRING2": "second string", "APP_SECTION2_VARBOOL2": "true"} {
			require.NoError(t, os.Setenv(k, v))
		}
		var c TestConf
		i := New(&c)
		err := i.SetFromEnv("APP")
		require.Error(t, err)
		require.Equal(t, TestConf{}, c)
	})
}

func TestSetFromDBPositive(t *testing.T) {

	// Если в базе есть все переменные конфига, метод вернет заполненный конфиг и nil.
	t.Run("Successful reading from DB", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"key", "value"})
		rows.AddRow("SECTION1.VARINT1", "11")
		rows.AddRow("SECTION1.VARSTRING1", "first string")
		rows.AddRow("SECTION1.VARBOOL1", "true")
		rows.AddRow("SECTION2.VARINT2", "22")
		rows.AddRow("SECTION2.VARSTRING2", "second string")
		rows.AddRow("SECTION2.VARBOOL2", "true")

		mock.ExpectQuery("SELECT key, value FROM").WithArgs("config").WillReturnRows(rows)
		var c TestConf
		i := New(&c)
		err := i.SetFromDB(db, "config")
		require.NoError(t, err)
		require.Equal(t, 11, c.Section1.VarInt1)
		require.Equal(t, "first string", c.Section1.VarString1)
		require.Equal(t, true, c.Section1.VarBool1)
		require.Equal(t, 22, c.Section2.VarInt2)
		require.Equal(t, "second string", c.Section2.VarString2)
		require.Equal(t, true, c.Section2.VarBool2)
	})

	// Если в базе нет некоторых переменных, метод вернет конфиг заполненный насколько возможно и nil.
	t.Run("Partial reading from DB", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"key", "value"})
		rows.AddRow("SECTION2.VARINT2", "22")
		rows.AddRow("SECTION2.VARSTRING2", "second string")
		rows.AddRow("SECTION2.VARBOOL2", "true")

		mock.ExpectQuery("SELECT key, value FROM").WithArgs("config").WillReturnRows(rows)
		var c TestConf
		i := New(&c)
		err := i.SetFromDB(db, "config")
		require.NoError(t, err)
		require.Equal(t, 0, c.Section1.VarInt1)
		require.Equal(t, "", c.Section1.VarString1)
		require.Equal(t, false, c.Section1.VarBool1)
		require.Equal(t, 22, c.Section2.VarInt2)
		require.Equal(t, "second string", c.Section2.VarString2)
		require.Equal(t, true, c.Section2.VarBool2)
	})

	// Если в базе нет переменных, метод вернет zerovalue конфиг и nil.
	t.Run("No vars in DB", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"key", "value"})

		mock.ExpectQuery("SELECT key, value FROM").WithArgs("config").WillReturnRows(rows)
		var c TestConf
		i := New(&c)
		err := i.SetFromDB(db, "config")
		require.NoError(t, err)
		require.Equal(t, TestConf{}, c)
	})
}

func TestSetFromDBNegative(t *testing.T) {

	// Если некоторые типы перепутаны, метод вернет zerovalue конфиг и ошибку.
	t.Run("Unexpected types", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		rows := sqlmock.NewRows([]string{"key", "value"})
		rows.AddRow("SECTION1.VARINT1", "first string")
		rows.AddRow("SECTION1.VARSTRING1", "true")
		rows.AddRow("SECTION1.VARBOOL1", "affwefefasdasfdsda")
		rows.AddRow("SECTION2.VARINT2", "22")
		rows.AddRow("SECTION2.VARSTRING2", "second string")
		rows.AddRow("SECTION2.VARBOOL2", "true")

		mock.ExpectQuery("SELECT key, value FROM").WithArgs("config").WillReturnRows(rows)
		var c TestConf
		i := New(&c)
		err := i.SetFromDB(db, "config")
		require.Error(t, err)
		require.Equal(t, TestConf{}, c)
	})
}

func TestCombinePositive(t *testing.T) {
	// Успешное чтение из файла, окружения и базы
	// Считанные из базы переписывают считанные из окружения, которые переписывают считанные из файла
}

func TestCombineNegative(t *testing.T) {
	// Если файла не существует, метод вернет пустой конфиг и ошибку.
	// Если файл не является TOML-readable, метод вернет пустой конфиг и ошибку.
	// Если некоторые типы перепутаны, метод вернет zerovalue конфиг и ошибку.
}

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
