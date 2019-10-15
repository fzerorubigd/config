package config_test

import (
	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Config_Set_Get_With_A_Simple_Key_String_Value(t *testing.T) {
	c, err := config.New()
	assert.NoError(t, err)

	c.Set("k", "v")
	v, err := c.Get("k")

	assert.NoError(t, err)
	assert.Equal(t, "v", v)
}

func Test_Config_Feed_With_Map_Repo(t *testing.T) {
	c, err := config.New(config.Options{Feeder: feeder.Map{
		"name":     "Hey You",
		"band":     "Pink Floyd",
		"year":     1979,
		"duration": 4.6,
	}})
	assert.NoError(t, err)

	v, err := c.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "Hey You", v)

	v, err = c.GetString("name")
	assert.NoError(t, err)
	assert.Equal(t, "Hey You", v)

	band, err := c.Get("band")
	assert.NoError(t, err)
	assert.Equal(t, "Pink Floyd", band)

	year, err := c.Get("year")
	assert.NoError(t, err)
	assert.Equal(t, 1979, year)

	year, err = c.GetInt("year")
	assert.NoError(t, err)
	assert.Equal(t, 1979, year)

	duration, err := c.Get("duration")
	assert.NoError(t, err)
	assert.Equal(t, 4.6, duration)

	duration, err = c.GetFloat("duration")
	assert.NoError(t, err)
	assert.Equal(t, 4.6, duration)

	wrong, err := c.Get("wrong.nested")
	assert.Error(t, err)
	assert.Equal(t, nil, wrong)
}

func Test_Config_Feed_With_Map_Repo_Includes_A_Slice(t *testing.T) {
	c, err := config.New(config.Options{Feeder: feeder.Map{
		"scores": map[string]interface{}{
			"A": 1,
			"B": 2,
			"C": 3,
		},
	}})
	assert.NoError(t, err)

	v, err := c.Get("scores.A")
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = c.Get("scores.B")
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
}

func Test_Config_Feed_It_Should_Get_Env_From_OS(t *testing.T) {
	err := os.Setenv("URL", "https://miladrahimi.com")
	if err != nil {
		panic(err)
	}

	c, err := config.New(config.Options{Feeder: feeder.Map{
		"url": "${ URL }",
	}})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)

	assert.Equal(t, os.Getenv("URL"), v)
}

func Test_Config_Feed_It_Should_Get_Env_Default_When_Not_In_OS(t *testing.T) {
	err := os.Setenv("EMPTY", "")
	if err != nil {
		panic(err)
	}

	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"url": "${ EMPTY | http://localhost }",
		},
	})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost", v)
}

func Test_Config_Feed_JSON(t *testing.T) {
	c, err := config.New(config.Options{
		Feeder: feeder.Json{Path: "test/config.json"},
	})
	assert.NoError(t, err)

	v, err := c.Get("numbers.2")
	assert.NoError(t, err)
	assert.Equal(t, float64(3), v)

	v, err = c.Get("users.0.address.city")
	assert.NoError(t, err)
	assert.Equal(t, "Delfan", v)
}
