const evtSource = new EventSource("/updates")
const esGoRoutines = new EventSource("/goroutines")

activeGRs = [];
Sites = new Map([])
rColors = []

var chartSate = "bar";

esGoRoutines.onmessage = (event) =>{
    document.getElementById("actives").innerText = "Active goroutines:"+String(event.data)
    newGRcount = parseInt(event.data)
    activeGRs.push(newGRcount)
    if (activeGRs.length > 70){
        activeGRs.shift()
    }
    grchart.update()
}

evtSource.onmessage = (event) => {
    console.log(String(event.data))
    document.getElementById("status").innerText = "Connected ✔"
    document.getElementById("status").style.backgroundColor = "#17fc03"
    document.getElementById("status").style.borderColor = "#12cc02"
    if(String(event.data).slice(0,1)=="b"){
        if(String(event.data).slice(1,)=="true"){
            document.getElementById("curMode").innerText = "Blocking mode"
        }else{
            document.getElementById("curMode").innerText = "Forwarding mode"
        }
    }else if (String(event.data).slice(0,2) == "co"){
        var c = String(event.data).slice(3)
        if (Sites.get(c) == undefined){
            Sites.set(c,1)
        }else{
            Sites.set(c,Sites.get(c)+1)
        }
        rColors.push("#"+String(Math.floor(Math.random()*16777215).toString(16)))
        connsChart.data.labels = Array.from(Sites.keys())
        connsChart.data.datasets[0].data = Array.from(Sites.values())
        connsChart.data.datasets[0].backgroundColor = rColors
        connsChart.update()
    }else {
        const pElement = document.createElement('p')
        pElement.className = "log"
        
        pElement.innerText = new Date().toLocaleDateString() + " "+ new Date().toLocaleTimeString() +" "+ String(event.data)
        var logs = document.getElementById("logs")
        logs.appendChild(pElement)
        for(var i = 0; i < logs.getElementsByClassName("log").length - parseInt(document.getElementById("storedlogs").value); i++){
            logs.removeChild(logs.getElementsByClassName("log").item(0))
        }
    }
}

evtSource.onerror = (event) =>{
    // handle dropped or failed connection
    console.log(event)
    document.getElementById("status").innerText = "Waiting for connection 🛑"
        document.getElementById("status").style.backgroundColor = "#b30419"
        document.getElementById("status").style.borderColor = "#960012"
  }

async function postData(url = '', data = '') {
    const resp = fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "text/plain"
        },
        body: data
    })
    return (await resp).text();
}

function stopServer(){
    postData("/reqs?q=stop")
}


function changeDisplay(){
    var btn = document.getElementById("cd")
    if (chartSate == "bar"){
        chartSate = "doughnut"
        btn.innerText = "Doughnut Chart"
        connsChart.destroy()
        connsChart = new Chart(connsCtx, {
            type: "doughnut",
            data: {
                labels: Array.from(Sites.keys()),
                datasets: [{
                    data: Array.from(Sites.values()),
                    backgroundColor: rColors,
                }],
            },
        })
    }else if (chartSate == "doughnut"){
        chartSate = "bar"
        btn.innerText = "Bar Chart"
        connsChart.destroy()
        connsChart = new Chart(connsCtx, {
            type: "bar",
            data: {
                labels: Array.from(Sites.keys()),
                datasets: [{
                    data: Array.from(Sites.values()),
                    backgroundColor: rColors,
                }],
            },
        })
    }
    connsChart.update()
}
