import React, {Component} from 'react';
import './SBox.scss'

class SBox extends Component {

    render() {
        return (
            <div className="weather34box">
                <div className="title">
                    <svg id="i-info" viewBox="0 0 32 32" width="10" height="10" fill="none" stroke="currentcolor"
                         stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%">
                        <path d="M16 14 L16 23 M16 8 L16 10"></path>
                        <circle cx="16" cy="16" r="14"></circle>
                    </svg>
                    {this.props.title}
                </div>
                {this.props.children}
            </div>
        );
    }
};

export default SBox
