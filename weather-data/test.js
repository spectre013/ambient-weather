

const fields = ['tempf','tempinf','temp1f','temp2f','baromrelin','uv','humidity','windspeedmph','windgustmph','dewpoint','humidityin','humidity1','humidity2','dailyrainin','lightiningday'];

fields.forEach((field) => {
    console.log("CREATE INDEX "+field+"IDX ON records ("+field+", date);");
})