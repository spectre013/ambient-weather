import React from "react";
import { useOutletContext } from "react-router-dom";

// Editorial about-the-station page. Speaks in the voice of a backyard
// observer / amateur climatologist. Static prose, two columns of running text,
// and a small specs table — no decorative imagery.

export default function AboutPage() {
  const { forecast } = useOutletContext();
  const lat = forecast.latitude.toFixed(4);
  const lon = forecast.longitude.toFixed(4);

  return (
    <article className="page about-page">
      <header className="page-head">
        <div className="page-meta">§ Vol 4 · Colophon</div>
        <h1 className="page-title">About Lorson Ranch.</h1>
        <p className="page-dek">
          Lorson Ranch Weather is a personal weather station running an Ambient
          WS-5000 array on my roof in Lorson Ranch, Colorado. It reports
          every sixteen seconds and the data being relayed to the website every
          minute.
        </p>
      </header>

      <div className="prose-grid">
        <div className="prose-col">
          <p>
            <span className="dropcap">T</span>he station occupies a quiet spot
            at five thousand seven hundred and thirty feet above sea level, on
            the south slope of a residential block that catches Pikes Peak
            wind from the west and the long, dry afternoons of the Front Range
            sun. It sees more weather than its size suggests.
          </p>
          <p>
            Readings are pushed from the indoor console to a service that translates the data
            to json and then pushes the json to a messaging server. The processing service
            retrieves that messages sand does some calculations and than persists the data
            to a Postgres instance running on the same machine. A Visual Crossing forecast
            is fetched every four hours and cached locally. All data is captured and stored
            locally with the only data coming in from external connections is forecast data.
          </p>
          <p>
            The dashboard you are reading was written in React. It draws its visual language
            from the printed page — mastheads, sub-rules, hairline columns, monospace captions — and
            it has been deliberately under-decorated. The intent is the
            instrument, not the chrome around it.
          </p>
        </div>

        <div className="prose-col">
          <p>
            <span className="dropcap">D</span>ata in <em>Climate</em> covers
            seven calendar years of monthly aggregates. The station was
            installed in late July 2020, so the earliest months are
            necessarily blank. Anything dated 2026 is the current year-to-date
            and will continue to fill in.
          </p>
          <p>
            I am not a meteorologist. I am an observer with a database. The pages
            here are reading instruments — they tell you what was measured,
            with as little editorialising. Where the underlying
            number is suspect, the cell is shown faint; where it is missing, it is empty.
          </p>
          <p>
            Source code lives at <code>github.com/spectre013/ambient-weather</code>
            and is offered under MIT. Corrections, calibration notes, and
            anomalous-reading reports are welcome by email.
          </p>
        </div>
      </div>

      <div className="page-rule" />

      <section className="specs">
        <h2 className="specs-title">§ Specifications</h2>
        <dl className="specs-grid">
          <div><dt>Station</dt><dd>Ambient WS-5000</dd></div>
          <div><dt>Console</dt><dd>Ambient Weather · firmware v4.3.7</dd></div>
          <div><dt>Reporting</dt><dd>16-second interval · WebSocket</dd></div>
          <div><dt>Backbone</dt><dd>Multiple Homelab servers · Postgres 16</dd></div>
          <div><dt>Forecast</dt><dd>Visual Crossing · NOAA NDFD</dd></div>
          <div><dt>Position</dt><dd>{lat}°N · {Math.abs(lon).toFixed(4)}°W</dd></div>
          <div><dt>Elevation</dt><dd>5,730 ft / 1,746 m</dd></div>
          <div><dt>Sensors</dt><dd>Anemometer · rain · UV · AQI · lightning · 5× remote</dd></div>
        </dl>
      </section>
    </article>
  );
}
