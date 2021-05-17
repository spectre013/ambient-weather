import moment from "moment";

const year = function(date) {
    return moment(date).format('YYYY');
}
const month = function(date) {
    return moment(date).format('MMM');
}

export {year,month}
