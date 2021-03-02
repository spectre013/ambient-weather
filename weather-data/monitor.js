#!/usr/bin/env node

const pm2 = require('pm2');

pm2.connect(function(err) {
  if (err) throw err;

    setTimeout(function worker() {
        console.log("Restarting data ...");
        pm2.restart('data', function() {});
        setTimeout(worker, 7200000);
        }, 7200000);
});