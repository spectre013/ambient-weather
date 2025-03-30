function alertColor(event) {
    /**
     * Convert alert event to CSS class name.
     * @param {string} event - The alert event name
     * @returns {string} CSS class name derived from the event
     */
    if (event && event.startsWith('911')) {
        return 'telephone-outage-911';
    }
    return event ? event.toLowerCase().replace(/\s+/g, '-') : '';
}

