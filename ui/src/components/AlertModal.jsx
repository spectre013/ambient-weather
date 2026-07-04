import React, { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { slugifyEvent } from "../utils/format";

const fmtDateTime = (v) =>
    v
        ? new Date(v).toLocaleString("en-US", {
              weekday: "short",
              month: "short",
              day: "numeric",
              hour: "numeric",
              minute: "2-digit",
          })
        : "—";

// Full-screen modal that shows the complete NWS advisory text. When several
// alerts are active it renders a sidebar list so the user can read every one.
export default function AlertModal({ alerts, index = 0, onClose }) {
    const list = (alerts || []).filter(Boolean);
    const [sel, setSel] = useState(index);

    // Keep selection in range and reset when opened at a new index.
    useEffect(() => {
        setSel(index < list.length ? index : 0);
    }, [index, list.length]);

    // Close on Escape and lock background scroll while open.
    useEffect(() => {
        const onKey = (e) => {
            if (e.key === "Escape") onClose();
        };
        window.addEventListener("keydown", onKey);
        const prev = document.body.style.overflow;
        document.body.style.overflow = "hidden";
        return () => {
            window.removeEventListener("keydown", onKey);
            document.body.style.overflow = prev;
        };
    }, [onClose]);

    if (!list.length) return null;
    const a = list[sel] || list[0];
    const cls = slugifyEvent(a.event);
    const multi = list.length > 1;

    return createPortal(
        <div className="modal-backdrop" onClick={onClose}>
            <div
                className={`modal nws-${cls}`}
                role="dialog"
                aria-modal="true"
                aria-label={a.event}
                onClick={(e) => e.stopPropagation()}
            >
                <div className="modal-head">
                    <div className="modal-head-main">
                        <div className="sev">{a.severity}</div>
                        <h2 className="modal-title">{a.event}</h2>
                        <div className="modal-sub">{a.senderName}</div>
                    </div>
                    <button className="modal-close" onClick={onClose} aria-label="Close">
                        ×
                    </button>
                </div>

                <div className="modal-cols">
                    {multi && (
                        <nav className="modal-list" aria-label="Active alerts">
                            <div className="modal-list-head">
                                {list.length} active alerts
                            </div>
                            {list.map((it, i) => (
                                <button
                                    key={it.id ?? i}
                                    className={`modal-list-item nws-${slugifyEvent(
                                        it.event
                                    )} ${i === sel ? "is-active" : ""}`}
                                    onClick={() => setSel(i)}
                                >
                                    <span className="mli-event">{it.event}</span>
                                    <span className="mli-meta">
                                        {it.severity} · until {fmtDateTime(it.end || it.expires)}
                                    </span>
                                </button>
                            ))}
                        </nav>
                    )}

                    <div className="modal-detail">
                        {a.headline && <p className="modal-headline">{a.headline}</p>}

                        <dl className="modal-facts">
                            <div>
                                <dt>Effective</dt>
                                <dd>{fmtDateTime(a.effective || a.sent)}</dd>
                            </div>
                            {a.onset && (
                                <div>
                                    <dt>Onset</dt>
                                    <dd>{fmtDateTime(a.onset)}</dd>
                                </div>
                            )}
                            <div>
                                <dt>Ends</dt>
                                <dd>{fmtDateTime(a.end || a.expires)}</dd>
                            </div>
                            <div>
                                <dt>Urgency</dt>
                                <dd>{a.urgency}</dd>
                            </div>
                            <div>
                                <dt>Certainty</dt>
                                <dd>{a.certainty}</dd>
                            </div>
                            {a.response && (
                                <div>
                                    <dt>Response</dt>
                                    <dd>{a.response}</dd>
                                </div>
                            )}
                        </dl>

                        {a.areadesc && (
                            <section className="modal-section">
                                <h3>Affected area</h3>
                                <p>{a.areadesc}</p>
                            </section>
                        )}

                        {a.description && (
                            <section className="modal-section">
                                <h3>Details</h3>
                                <p className="pre">{a.description}</p>
                            </section>
                        )}

                        {a.instruction && (
                            <section className="modal-section">
                                <h3>Instructions</h3>
                                <p className="pre">{a.instruction}</p>
                            </section>
                        )}
                    </div>
                </div>
            </div>
        </div>,
        document.body
    );
}
