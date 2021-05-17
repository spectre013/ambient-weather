import React,{Component} from "react";
import './RainSmall.scss'
import {year,month} from '../common/Dates'

class RainSmall extends Component {



    render() {
        console.log(this.props)
        return (
            <div className="value">
                <div id="position3">
                    <div className="topmin">
                        <div className="topblue1 bluebg">{this.props.conditions.monthlyrainin }
                            <span className="smallwindunit">in </span>
                        </div>
                    </div>
                    <div className="minword">{ month(this.props.conditions.date) }</div>
                    <div className="mintimedate">Total
                    </div>
                    <div className="yearwordbig">{ year(this.props.conditions.date) }</div>
                    <div className="topmax">
                        <div className="topblue1 bluebg">{ this.props.conditions.yearlyrainin }
                            <span className="smallwindunit">in </span>
                        </div>
                    </div>
                    <div className="maxword">{ year(this.props.conditions.date) }</div>
                    <div className="maxtimedate">Total</div>
                </div>
            </div>
        )
    }
};

export default RainSmall;
