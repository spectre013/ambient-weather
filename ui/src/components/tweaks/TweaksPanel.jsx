import React, { useEffect, useRef, useState } from "react";

// Tweaks panel — collapsible bottom-right control surface, with the host
// postMessage protocol so the toolbar toggle works.
//
// useTweaks(defaults) returns [state, setKey] where setKey can be called as
//   setKey("foo", value)         (single key)
//   setKey({ foo, bar })         (merge object)

export function useTweaks(defaults) {
  const [state, setState] = useState(defaults);
  const set = (a, b) => {
    let edits;
    if (typeof a === "string") edits = { [a]: b };
    else edits = a;
    setState((s) => ({ ...s, ...edits }));
    try {
      window.parent.postMessage({ type: "__edit_mode_set_keys", edits }, "*");
    } catch {}
  };
  return [state, set];
}

export function TweaksPanel({ title = "Tweaks", children }) {
  const [open, setOpen] = useState(false);
  const ready = useRef(false);

  useEffect(() => {
    const onMsg = (e) => {
      const d = e.data || {};
      if (d.type === "__activate_edit_mode") setOpen(true);
      else if (d.type === "__deactivate_edit_mode") setOpen(false);
    };
    window.addEventListener("message", onMsg);
    if (!ready.current) {
      ready.current = true;
      try { window.parent.postMessage({ type: "__edit_mode_available" }, "*"); } catch {}
    }
    return () => window.removeEventListener("message", onMsg);
  }, []);

  const close = () => {
    setOpen(false);
    try { window.parent.postMessage({ type: "__edit_mode_dismissed" }, "*"); } catch {}
  };

  if (!open) return null;
  return (
    <div className="tweaks-panel">
      <header>
        <span>{title}</span>
        <button onClick={close} aria-label="Close">×</button>
      </header>
      <div className="body">{children}</div>
    </div>
  );
}

export const TweakSection = ({ label, children }) => (
  <div className="tweak-section">
    <div className="tweak-section-label">{label}</div>
    {children}
  </div>
);

export const TweakRadio = ({ label, value, options, onChange }) => (
  <div className="tweak-row">
    <span className="tweak-label">{label}</span>
    <div className="tweak-radio">
      {options.map((o) => (
        <button key={o.value}
          className={value === o.value ? "active" : ""}
          onClick={() => onChange(o.value)}>{o.label}</button>
      ))}
    </div>
  </div>
);

export const TweakColor = ({ label, value, onChange, options }) => (
  <div className="tweak-row">
    <span className="tweak-label">{label}</span>
    <div className="tweak-swatches">
      {options.map((c) => (
        <button key={c}
          className={`swatch ${value === c ? "active" : ""}`}
          onClick={() => onChange(c)}
          style={{ background: c }}
          aria-label={c} />
      ))}
    </div>
  </div>
);

export const TweakToggle = ({ label, value, onChange }) => (
  <div className="tweak-row">
    <span className="tweak-label">{label}</span>
    <button
      className={`tweak-toggle ${value ? "on" : ""}`}
      onClick={() => onChange(!value)}>
      <span />
    </button>
  </div>
);
