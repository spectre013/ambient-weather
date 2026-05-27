import React from "react";
import { NavLink } from "react-router-dom";

const links = [
  { to: "/",        label: "Dashboard", index: "01" },
  { to: "/climate", label: "Climate",   index: "02" },
  { to: "/about",   label: "About",     index: "03" },
];


function SegToggle({ value, options, onChange, ariaLabel }) {
    return (
        <div className="nav-seg" role="group" aria-label={ariaLabel}>
            {options.map((o) => (
                <button
                    key={o.value}
                    type="button"
                    className={value === o.value ? "active" : ""}
                    onClick={() => onChange(o.value)}
                    aria-pressed={value === o.value}
                >
                    {o.label}
                </button>
            ))}
        </div>
    );
}

export default function Nav({ units, onUnitsChange, theme, onThemeChange }) {
    return (
        <nav className="primary-nav">
            <ul>
                {links.map((l) => (
                    <li key={l.to}>
                        <NavLink to={l.to} end={l.to === "/"}>
                            {({ isActive }) => (
                                <>
                                    <span className="ix">{l.index}</span>
                                    <span className="lb">{l.label}</span>
                                    {isActive && <span className="rule" />}
                                </>
                            )}
                        </NavLink>
                    </li>
                ))}
            </ul>
            <div className="nav-controls">
                <SegToggle
                    ariaLabel="Units"
                    value={units}
                    onChange={onUnitsChange}
                    options={[{ value: "F", label: "°F" }, { value: "C", label: "°C" }]}
                />
                <SegToggle
                    ariaLabel="Theme"
                    value={theme}
                    onChange={onThemeChange}
                    options={[
                        { value: "light", label: (<><SunGlyph /></>) },
                        { value: "dark",  label: (<><MoonGlyph /></>) },
                    ]}
                />
            </div>
        </nav>
    );
}

function SunGlyph() {
    return (
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none"
             stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2M12 20v2M2 12h2M20 12h2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
        </svg>
    );
}
function MoonGlyph() {
    return (
        <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" />
        </svg>
    );
}
