import Vue from 'vue';
import App from './App.vue';
import computed from './mixins/temperature';

Vue.config.productionTip = false
Vue.config.ignoredElements = [/^darksky/,/^temp/,/^weather34/,'oorange','span2','span4','orange','ogreen','spanelightning','ored','heatindex','spanmaxwind',
'rainf','spanwindtitle','spanmaxwind','oblue','yesterdaytimemax','spaneindoortemp','suptemp','trendmovementrising','trendmovementfalling','rainblue','grey',
'spancalm','supmb','max','supmb','windrun','darkgrey','convtext','steady','strongnumbers','period','min','minutes','hrs','spanindoortempfalling',
'indoororange1','spanindoortempfalling','indoororange1','spanindoortempfalling','luminance1','uppercase','span1','uviforecasthourgreen',
'value','valuetext','smalluvunit','alertvalue','noalert','spanyellow','valuetitleunit','topblue1','smallwindunit','smallrainunita','valuetextheading1','raiblue',
  'smallrainunit2','rainratetextheading','maxred','hours','blueu','moonrisecolor','moonm','moonsetcolor','tgreen','topgreen1','smallwindunit','toporange1','minblue',
  'smalltempunit','smalltempunit2','smallrainunit','valuetext1','uviforecasthouryellow','trendmovementfallingx','redu','blue','articlegraph','uviforecasthourred'];


Vue.mixin(computed);

new Vue({
  render: h => h(App),
}).$mount('#app')
