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
        <h1 className="page-title">A small instrument in a small yard.</h1>
        <p className="page-dek">
          The Backyard Observer is a personal weather station running an Ambient
          WS-2902 array on a back fence in Lorson Ranch, Colorado. It reports
          every sixteen seconds and has not missed a calibration in fourteen
          months.
        </p>
      </header>

      <div className="prose-grid">
        <div className="prose-col">
          <p>
            <span className="dropcap">T</span>he station occupies a quiet spot
            at five thousand one hundred and eighty feet above sea level, on
            the south slope of a residential block that catches Pikes Peak
            wind from the west and the long, dry afternoons of the Front Range
            sun. It sees more weather than its size suggests.
          </p>
          <p>
            Readings are pushed from the indoor console to a Raspberry Pi 4,
            collected by a small Go program, and persisted to a Postgres
            instance running on the same machine. A Visual Crossing forecast
            is fetched once an hour and cached locally. There is no cloud, no
            account, no app store; just a coaxial cable, a hex wrench, and a
            text file of station offsets that I keep updating each spring.
          </p>
          <p>
            The dashboard you are reading was written in React over a long,
            cold January. It draws its visual language from the printed page —
            mastheads, sub-rules, hairline columns, monospace captions — and
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
            I am not a meteorologist. I am an observer with a soldering iron
            and a database. The pages here are reading instruments — they
            tell you what was measured, with as little editorialising as I
            could manage. Where the underlying number is suspect, the cell is
            shown faint; where it is missing, it is empty.
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
          <div><dt>Station</dt><dd>Ambient WS-2902 (osprey-mount)</dd></div>
          <div><dt>Console</dt><dd>WS-2902-C · firmware v4.3.5</dd></div>
          <div><dt>Reporting</dt><dd>16-second interval · WebSocket</dd></div>
          <div><dt>Backbone</dt><dd>Raspberry Pi 4 · Postgres 16</dd></div>
          <div><dt>Forecast</dt><dd>Visual Crossing · NOAA NDFD</dd></div>
          <div><dt>Position</dt><dd>{lat}°N · {Math.abs(lon).toFixed(4)}°W</dd></div>
          <div><dt>Elevation</dt><dd>5,180 ft / 1,579 m</dd></div>
          <div><dt>Sensors</dt><dd>Anemometer · rain · UV · AQI · lightning · 5× remote</dd></div>
          <div><dt>Last calibration</dt><dd>April 4, 2026</dd></div>
          <div><dt>Uptime</dt><dd>14 d (rolling)</dd></div>
        </dl>
      </section>
    </article>
  );
}
