<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body">
                    <div class="weather34darkbrowser" url="Forecast Summary"></div>
                    <main class="grid">
                        <article>
                            <h1>Hourly Summary</h1>
                        </article>
                    </main>
                </div>
            </section>
            <button class="lity-close" type="button" aria-label="Close (Press escape to close)" @click="close">×</button>
        </div>
    </div>
</template>


<script>
    export default {
        name: 'forecastsummaryymodel',
        props: {
            forecast: Object
        },
        methods: {
            close: function() {
                this.$parent.closeModal('forecastsummary');
            }
        },
        mounted() {

        },
        computed: {
        },
    };
</script>

<style scoped>

</style>