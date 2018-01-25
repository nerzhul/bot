package internal

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
)

func (r *commandRouter) handlerWeather(args string, user string, channel string) *string {
	w, err := owm.NewCurrent(gconfig.OpenWeatherMap.Unit, gconfig.OpenWeatherMap.Lang, gconfig.OpenWeatherMap.APIKey)
	if err != nil {
		log.Errorf("Unable to initialize OpenWeatherMap: %s", err.Error())
		return nil
	}

	err = w.CurrentByName(args)
	if err != nil {
		log.Errorf("Unable to fetch OpenWeatherMap data: %s", err.Error())
		return nil
	}

	result := new(string)
	if len(w.Weather) == 0 {
		*result = fmt.Sprintf("Weather for %s: no desc | Wind speed: %.2f | Humidity: %d%% | "+
			"Temperature %.2f° (min %.2f°, max %.2f°)", args, w.Wind.Speed, w.Main.Humidity,
			w.Main.Temp, w.Main.TempMin, w.Main.TempMax)
	} else {
		*result = fmt.Sprintf("Weather for %s: %s | Wind speed: %.2fm/s | Humidity: %d%% | "+
			"Temperature %.2f° (min %.2f°, max %.2f°)", args, w.Weather[0].Description, w.Wind.Speed, w.Main.Humidity,
			w.Main.Temp, w.Main.TempMin, w.Main.TempMax)
	}

	return result
}
