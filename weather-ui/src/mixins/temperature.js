export default {
    filters: {
        dewPointClass: function(dewpoint) {
            if (dewpoint>69.8) {
                return 'tempmodulehome25-30c';
            } else if (dewpoint>=68) {
                return 'tempmodulehome20-25c';
            } else if (dewpoint>=59) {
                return 'tempmodulehome15-20c';
            } else if (dewpoint>=50) {
                return 'tempmodulehome10-15c';
            } else if(dewpoint>41){
                return 'tempmodulehome5-10c';
            } else if(dewpoint>=32){
                return 'tempmodulehome0-5c';
            } else if(dewpoint>14){
                return 'tempmodulehome-10-0c';
            } else if(dewpoint>=-50){
                return 'tempmodulehome-50-10c';
            }
        },
        humidityClass: function(humidity) {
            if (humidity>90){
                return 'temphumcircle80-100'
            } else if (humidity>70){
                return 'temphumcircle60-80'
            } else if (humidity>35){
                return 'temphumcircle35-60'
            } else if (humidity>25){
                return 'temphumcircle25-35'
            } else if (humidity<=25){
                return 'temphumcircle0-25'
            }
            return "";
        },
        temperaturetoday: function(temp) {
            if (temp>=105.8){
                return "temperaturetoday41-45"
            } else if (temp>=96.8){
                return "temperaturetoday36-40"
            } else if (temp>=87.8) {
                return "temperaturetoday31-35"
            } else if (temp>=78.8) {
                return "temperaturetoday26-30"
            } else if (temp>=69.8) {
                return "temperaturetoday21-25"
            } else if (temp>=60.8) {
                return "temperaturetoday16-20"
            } else if (temp>=50) {
                return "temperaturetoday11-15"
            } else if (temp>42.8) {
                return "temperaturetoday6-10"
            } else if (temp>=32) {
                return "temperaturetoday0-5"
            } else if (temp<32) {
                return "temperaturetodayminus"
            } else if (temp<=23) {
                return "temperaturetodayminus5"
            } else if (temp<-14) {
                return "temperaturetodayminus10"
            }
            return "";
        },
        smallTempClass: function(temp) {
            if (temp>=104){
                return "tempmodulehome40-50c"
            } else if (temp>=95){
                return "tempmodulehome35-40c"
            } else if (temp>=86) {
                return "tempmodulehome30-35c"
            } else if (temp>=77) {
                return "tempmodulehome25-30c"
            } else if (temp>=68) {
                return "tempmodulehome20-25c"
            } else if (temp>=59) {
                return "tempmodulehome15-20c"
            } else if (temp>=50) {
                return "tempmodulehome10-15c"
            } else if (temp>41) {
                return "tempmodulehome5-10c"
            } else if (temp>=32) {
                return "tempmodulehome0-5c"
            } else if (temp>14) {
                return "tempmodulehome-10-0c"
            } else if (temp>-50) {
                return "tempmodulehome-50-10c"
            }
            return "";
        },
        tempClass: function(temp) {
            if(temp < 14) {
                return "outsideminus10";
            } else if(temp <= 23) {
                return "outsideminus5";
            } else if(temp <= 32) {
                return "outsidezero";
            } else if(temp < 41) {
                return "outside0-5";
            } else if(temp < 50) {
                return "outside6-10";
            } else if(temp < 59) {
                return "outside11-15";
            } else if(temp < 68) {
                return "outside16-20";
            } else if(temp < 77) {
                return "outside21-25";
            } else if(temp < 86) {
                return "outside26-30";
            } else if(temp < 95) {
                return "outside31-35";
            } else if(temp < 104) {
                return " outside36-40";
            } else if(temp < 113) {
                return " outside41-45";
            } else if(temp < 150) {
                return " outside50";
            }
            return "";
        }
    }
}



