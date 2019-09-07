import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        settings:{}
    },
    mutations: {
        setTheme(state, theme) {
            state.settings.theme = theme;
        },
        setUnits(state, units) {
            state.settings.units = units;
        },
        setSpeed(state, speed) {
            state.settings.speed = speed;
        },
        setSettings(state,settings) {
            state.settings = settings;
        }
    },
    actions: {
        getSettings(context) {
            let settings = {'theme':'dark','units':'imperial','speed':'mph'}
            if (typeof(Storage) !== "undefined") {
                let storedSettings = window.localStorage.getItem('settings');
                if(!storedSettings) {
                    window.localStorage.setItem('settings',JSON.stringify(settings));
                    //console.log('default',settings);
                } else {
                    settings = JSON.parse(storedSettings);
                    //console.log('stored',settings);
                }
            }
            context.commit('setSettings', settings)
        },
        setSettings(context,settings) {
            if (typeof (Storage) !== "undefined") {
                window.localStorage.setItem('settings',JSON.stringify(settings));
            }
            context.commit('setSettings', settings)
        },
        setUnits(context, units) {
            context.commit('setUnits', units)
        },
        setTheme(context, theme) {
            context.commit('setTheme', theme)
        },
        setSpeed(context, speed) {
            context.commit('setSpeed', speed)
        }
    },
    getters: {
        settings: (state) => {
            return state.settings;
        },
        theme: (state) => {
            console.log("instore", state.settings.theme);
            return state.settings.theme;
        },
        units: (state) => {
            return state.settings.units;
        },
        speed: (state) => {
            return state.settings.speed;
        },
    }

})