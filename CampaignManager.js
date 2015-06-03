/**
 * Created by Pawan on 6/1/2015.
 */
var stringify=require('stringify');
var redis=require('redis');
var messageFormatter = require('DVP-Common/CommonMessageGenerator/ClientMessageJsonFormatter.js');
var DbConn = require('DVP-DBModels');
var config = require('config');

var port = config.Redis.port;
var ip = config.Redis.ip;
var hpath=config.Host.hostpath;
var logger = require('DVP-Common/LogHandler/CommonLogHandler.js').logger;

var client = redis.createClient("127.0.0.1",6379);
client.on("error", function (err) {
    console.log("Error " + err);

});

function AddCampaign(req,callback)
{
    // DbConn.
    try
    { var obj=req.body;
        var CampaignObject = DbConn.Campaign
            .build(
            {
                CampaignName:obj.CampaignName,
                CampaignNumber: obj.CampaignNumber,
                Max: obj.Max,
                Min:obj.Min,
                StartTime:obj.StartTime,
                EndTime: obj.EndTime,
                Enable: obj.Enable,
                Limit: obj.Limit,
                ScheduleId:obj.ScheduleId,
                Class: "OBJCLZ",
                Type: "OBJTYP",
                Category: "OBJCAT",
                CompanyId: 1,
                TenantId: 1
                // AddTime: new Date(2009, 10, 11),
                //  UpdateTime: new Date(2009, 10, 12),
                // CSDBCloudEndUserId: jobj.CSDBCloudEndUserId


            }
        )
    }
    catch(ex)
    {
        logger.error('[DVP-DialerApi.NewCampaign] - [%s] - [PGSQL] - Exception occurred while Saving Campaign Data ',req.body,ex);
        callback(ex, undefined);
    }

    CampaignObject.save().complete(function (err, result) {

        if(err)
        {
            logger.error('[DVP-DialerApi.NewCampaign] - [%s] - [PGSQL] - New Schedule %s saving failed unsuccessful',req.body,err);
            //var jsonString = messageFormatter.FormatMessage(err, "AppObject saving error", false, result);
            callback(err, undefined);
        }else{
            logger.debug('[DVP-DialerApi.NewCampaign] - [%s] - [PGSQL] - New Campaign %s is added successfully',req.body);
            callback(undefined,result);
        }


    });

}

function LoadCampaigns(req,callback)
{
    // DbConn.
    try
    {
        var obj=req.body;
        var CTime=moment().format("YYYY-MM-DD HH:mm");
        DbConn.Campaign.findAll({where: [{Enable: '1'}]}).complete(function (err,result)
        {
            if(err)
            {
                logger.error('[DVP-DialerApi.LoadCampaign] - [%s] - [PGSQL] - Exception occurred while Searching Campaign Data ',req.body,ex);
                callback(err, undefined);
            }
            else
            {
                if(result.length==0)
                {
                    logger.error('[DVP-DialerApi.LoadCampaign] - [%s] - [PGSQL] - No campaign found  ');
                    //  console.log('No user with the Extension has been found.');
                    ///logger.info( 'No user found for the requirement. ' );
                    callback('No Campaign found', undefined);
                }
                else
                {
                    for(var index in result)
                    {
                        if(CheckValidCampaign(result[index].StartTime.toString(),result[index].StartTime.toString()))
                        {
                            var CampName=result[index].CampaignName+"_"+result[index].id;
                            client.lpush("CMPLIST",CampName,function(err,reply)
                            {
                                if(err)
                                {
                                    logger.error('[DVP-DialerApi.LoadCampaign] - [%s] - [PGSQL] - Exception occurred while Pushing to redis',err);
//callback(err,undefined);
                                    continue;
                                }
                                else
                                {
                                    logger.debug('[DVP-DialerApi.LoadCampaign] - [%s] - [PGSQL] - Valid campaign picked');
                                    //callback(undefined,reply);
                                    SetCampaignMaxMin("Max",result[index].Max.toString(),CampName);
                                    SetCampaignMaxMin("Min",result[index].Min.toString(),CampName);
                                    FillCampaignPhones(result[index].id,CampName,result[index].Max.toString())
                                    if(index==result.length-1)
                                    {
                                        callback(undefined,"Done");
                                    }

                                }

                            })
                        }
                        else
                        {
                            continue;
                        }
                    }
                }

            }
        })


    }
    catch(ex)
    {
        logger.error('[DVP-DialerApi.LoadCampaign] - [%s] - [PGSQL] - Exception occurred while Loading Campaign Data ',req.body,ex);
        callback(ex, undefined);
    }



}

function PickCurrentCampaign(callback)
{
    client.lpop("CMPLIST",function(err,reply)
    {
        if(err)
        {
            callback(err,undefined)
        }
        else
        {
            callback(undefined,reply)
        }

    })
}

function GetPhonesOfCampaign(CampName)
{
    client.llen(CampName,function(err,result)
    {
        if(err)
        {

        }
        else
        {
            if(result>0)
            {
                client.lpop(CampName,function(errPop,resPop)
                {
                    if(errPop)
                    {

                    }
                    else
                    {
                        var count=CheckFillCount(CampName);
                        if(count>0)
                        {
                            var arr =CampName.split("_");
                            FillCampaignPhones(arr[1],CampName,count);
                            return resPop;
                        }
                    }
                })
            }
        }
    })
}


function CheckValidCampaign(StDt,EnDt)
{
    var x = moment(moment().format("YYYY-MM-DD HH:mm")).isBetween(StDt, EnDt);
    return x;
}

function SetCampaignMaxMin(MXMN,value,CampName)
{
    var CampMaxMin="";
    if(MXMN=="Max")
    {
        CampMaxMin=CampName+"_MAX";
    }else
    {
        CampMaxMin=CampName+"_MIN";
    }

    client.set(CampMaxMin,value,function(err,reply)
    {

    })



}

function FillCampaignPhones(campId,campName,Max)
{
    DbConn.CampaignPhones.findAll({where:[{CampaignId:campId},{Enable:"1"}]},{limit:Max}).complete(function(err,result)
    {
        if(err)
        {

        }
        else
        {

            for(var index in result)
            {
                client.lpush(campName,result[index].Phone.toString())
            }
        }
    })
}

function CheckFillCount(CampName)
{
    var max=CampName+"_MAX";
    var min=CampName+"_MIN";
    var MX=null;
    var MN=null;
    var LN=null;

    client.llen(CampName,function(errLen,resLen)
    {
        if(errLen)
        {

        }
        else
        {
            LN=resLen;
            MX=GetCurrentMaxMin(CampName,"MAX");
            MN=GetCurrentMaxMin(CampName,"MIN");

            if(LN<=MN)
            {
                return (MX-LN);
            }
            else
            {
                return 0;
            }
        }
    })

}

function GetCurrentMaxMin(CampName,MaxMin)
{
    var MxMinNm=CampName+"_"+MaxMin;
    client.get(MxMinNm,function(errMxMn,resMxMn)
    {
        if(errMin)
        {

        }
        else
        {
            return resMxMn;
        }
    })


}

module.exports.AddCampaign = AddCampaign;
module.exports.LoadCampaigns = LoadCampaigns;
module.exports.PickCurrentCampaign = PickCurrentCampaign;
module.exports.GetPhonesOfCampaign = GetPhonesOfCampaign;