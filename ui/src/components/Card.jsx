import React from "react";

// Generic card shell — single source of truth for the .card surface + header row.
export default function Card({ label, right, children, className = "", ...rest }) {
  return (
    <div className={`card ${className}`} {...rest}>
      <div className="label">
        <span>{label}</span>
        <span className="right">{right}</span>
      </div>
      {children}
    </div>
  );
}
