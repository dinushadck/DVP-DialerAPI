/**
 * Created by Pawan on 6/1/2015.
 */
var restify = require('restify');
var moment=require('moment');
var redis=require('redis');
var client = redis.createClient();
client.on("error", function (err) {
    console.log("Error " + err);

});
var RestServer = restify.createServer({
    name: "myapp",
    version: '1.0.0'
},function(req,res)
{

});
RestServer.listen(8081, function () {
    console.log('%s listening at %s', RestServer.name, RestServer.url);
   // var x = moment(moment().format("YYYY-MM-DD HH:mm")).isBetween("2013-12-09 12:31", "2016-12-09 16:31");
    //

    var p =client.llen("a",function(err,res){
        console.log(p);
    });

});

