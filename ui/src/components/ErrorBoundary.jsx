import React from "react";

// ErrorBoundary keeps one broken component from blanking the whole page.
// It shows a short notice in place of its children and tries them again
// whenever resetKey changes (e.g. on the next websocket update).
export default class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error) {
    return { error };
  }

  componentDidCatch(error, info) {
    console.error("render failed:", error, info?.componentStack);
  }

  componentDidUpdate(prev) {
    if (this.state.error && prev.resetKey !== this.props.resetKey) {
      this.setState({ error: null });
    }
  }

  render() {
    if (this.state.error) {
      return (
        <div className="card" role="alert" style={{ margin: "var(--gap) 0" }}>
          <div className="sub">Part of this page could not be displayed. It will retry on the next update.</div>
        </div>
      );
    }
    return this.props.children;
  }
}
