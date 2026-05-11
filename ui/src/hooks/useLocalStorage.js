// src/hooks/useLocalStorage.js
import { useState, useEffect } from "react";

export function useLocalStorage(key, initial) {
    const [value, setValue] = useState(() => {
        try {
            const raw = localStorage.getItem(key);
            return raw !== null ? JSON.parse(raw) : initial;
        } catch {
            return initial;
        }
    });

    useEffect(() => {
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch {
        }
    }, [key, value]);

    // Sync across tabs
    useEffect(() => {
        const onStorage = (e) => {
            if (e.key === key && e.newValue !== null) {
                try {
                    setValue(JSON.parse(e.newValue));
                } catch {
                }
            }
        };
        window.addEventListener("storage", onStorage);
        return () => window.removeEventListener("storage", onStorage);
    }, [key]);
    return [value, setValue];
}