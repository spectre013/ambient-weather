<template>
    <div class="modal-backdrop">
        <div class="modal">
            <div class="modal-header">
                <slot name="header">
                    <span class="close-icon" @click="close">X</span>
                </slot>
            </div>
            <section class="modal-body" v-if="forecast">
                <h1 v-bind:class="alertColor"><svg id="firealert" viewBox="0 0 32 32" width="22px" height="22px" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M16 3 L30 29 2 29 Z M16 11 L16 19 M16 23 L16 25"></path></svg>
                    {{ forecast.alerts[0].title}}</h1>
                <div class="alertText"><span v-bind:class="alertColor">Alert Description: </span> {{ forecast.alerts[0].description}}</div>
                <div class="targetarea"><span v-bind:class="alertColor">Regions Affected: </span> <span v-bind:key="region" v-for="region in forecast.alerts[0].regions">{{ region }}</span></div>
            </section>
        </div>
    </div>
</template>

<script>
  export default {
        name: 'alertmodal',
        props: {
            forecast: Object,
        },
        mounted: function() {
            console.log(this.forecast);
        },
        methods: {
            close: function() {
                this.$parent.closeModal('alert');
            },
        },
        computed: {
            alertColor: function() {
                let c = "";
                switch(this.forecast.alerts[0].title) {
                    case "Flash Flood Watch":
                    case "Flash Flood Warning":
                        c = 'greenu';
                        break;
                   case "Air Quality Alert":
                        c = 'orangeu';
                        break;
                    default:
                        c = "";
                }
                return c;
            }
        },
        watch: {

        },
    };
</script>

<style scoped>
    h1 {
        font-size: 22px;
        text-align: center;
    }
    .alertText {
        color: silver;
    }
    .targetarea {
        margin-top:10px;
        color: silver;
    }
    .modal-backdrop {
        position: fixed;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 10000;
    }

    .modal {
        background: #181818;
        box-shadow: 2px 2px 20px 1px;
        height: 550px;
        width: 900px;
        overflow-x: auto;
        display: flex;
        flex-direction: column;
    }
    .close-icon {
        position: relative;
        top: 0;
        left: 880px;
        cursor: pointer;
        text-align: center;
        font: 22px/25px arial,sans-serif;
        color: #9aba2f;
        font-weight: 900;
    }
    .modal-header {
        height: 25px;
        display: flex;
    }

    .modal-header {
        justify-content: space-between;
    }

    .modal-body {
        position: relative;
        padding: 20px 10px;
    }

    .btn-close {
        border: none;
        font-size: 20px;
        padding: 20px;
        cursor: pointer;
        font-weight: bold;
        color: #4AAE9B;
        background: transparent;
    }
    section {
        text-align: left;
    }
</style>