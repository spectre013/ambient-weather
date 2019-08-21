const pm2 = require('pm2');

pm2.connect(function(err) {
  if (err) throw err;

    setTimeout(function worker() {
        console.log("Restarting server ...");
         pm2.restart('server', function() {});
        setTimeout(worker, 7200000);
        }, 7200000);
});