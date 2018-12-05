let loadedProject = null;
let currentNodes = [];

function createProject() {
    let projName = prompt("New Project Name");
    if (projName) {
        var oReq = new XMLHttpRequest();
        oReq.open("GET", "/api/NewProject/" + projName);
        oReq.send();
    }
}

function addNode() {
    if (loadedProject) {
        let nodeName = prompt("New Node Name: ");
        if (nodeName) {
            var oReq = new XMLHttpRequest();
            console.log(loadedProject);
            oReq.open("GET", "/api/LoadProject/NewNode/" + nodeName + "/_/" + loadedProject);
            oReq.addEventListener("load", function(res){
                let tree = parseCallOutput(res.currentTarget.responseText, []);
                currentNodes = tree;
                drawNodes();
            });
            oReq.send();
        }
    } else {
        alert("You must load a project first");
    }
}

function link(a, b) {
    if (loadedProject) {
        if (a && b) {
            var oReq = new XMLHttpRequest();
            oReq.open("GET", "/api/LoadProject/Link/" + a + "/" + b + "/" + loadedProject);
            oReq.addEventListener("load", function(res){
                let tree = parseCallOutput(res.currentTarget.responseText, []);
                currentNodes = tree;
                drawNodes();
            });
            oReq.send();
        }
    } else {
        alert("You must load a project first");
    }
}

function linkNodes(obj) {
    $("#draw").attr("selected", "true");
    let clicked_a = false;
    let clicked_b = false;
    $("#draw").click(function(e) {
        console.log("Click");
        if (!clicked_a) {
            clicked_a = e.target;
        } else if (!clicked_b) {
            clicked_b = e.target;
            $("#draw").attr("selected", null);
            if (clicked_a != clicked_b) {
                if (clicked_a.getAttribute("attr-nodeId") && clicked_b.getAttribute("attr-nodeId")) {
                    link(clicked_a.innerText, clicked_b.innerText);
                }
            }
        }
    })
}

function removeNode(nodeName) {
    if (loadedProject) {
        if (nodeName) {
            var oReq = new XMLHttpRequest();
            console.log(loadedProject);
            oReq.open("GET", "/api/LoadProject/RemoveNode/" + nodeName + "/_/" + loadedProject);
            oReq.addEventListener("load", function(res){
                let tree = parseCallOutput(res.currentTarget.responseText, []);
                currentNodes = tree;
                drawNodes();
            });
            oReq.send();
        }
    } else {
        alert("You must load a project first");
    }
}

function loadProject(text) {
    var oReq = new XMLHttpRequest();
    oReq.open("GET", "/api/LoadProject/None/_/_/" + text);
    oReq.addEventListener("load", function(res){
        let tree = parseCallOutput(res.currentTarget.responseText, []);
        currentNodes = tree;
        drawNodes();
        loadedProject = text;
    });
    oReq.send();
}

function updateProjects() {
    var oReq = new XMLHttpRequest();
    oReq.addEventListener("load", function(res){
        let splits = res.currentTarget.responseText.split("\n");
        document.getElementById("projects").innerHTML = null;
        for (let i=0; i<splits.length; i++) {
            let node = document.createElement("button");
            node.innerText = splits[i];
            document.getElementById("projects").appendChild(node);
            node.setAttribute("onclick", "loadProject(this.innerText);");
        }
    });
    oReq.open("GET", "/api/GetProjects");
    oReq.send();
}

function parseCallOutput(str, nodeTree) {
    console.log(str);
    found = false;
    let splits = str.split("\n");
    let node;
    for (let i=0; i<splits.length; i++) {
        if (splits[i] == "<<GENERATE>>") {
            found = true;
        } else if (found) {
            let lineSplit = splits[i].split(":");
            if (lineSplit[0] == "Node") {
                node = {};
                node.id = Number(lineSplit[1]);
                node.value = lineSplit[2];
                nodeTree.push(node);
            } else if (lineSplit[0] == "Connection") {
                node.connections = node.connections || [];
                console.log("C:" + node.value + ":" + lineSplit[1]);
                node.connections.push(lineSplit[1]);
            }
        }
    }
    return nodeTree;
}

function drawNodes() {
    let connections = [];
    document.getElementById("draw").innerHTML = null;
    for (let i=0; i<currentNodes.length; i++) {
        let node = document.createElement("div");
        node.innerText = currentNodes[i].value;
        node.setAttribute("attr-nodeId", currentNodes[i].id);
        node.setAttribute("draggable", "true");
        node.className = "node";
        document.getElementById("draw").append(node);
        $(node).draggable({ 
            scroll: false,
            drag: function(){
                var offset = $(node).offset();
                var xPos = offset.left;
                var yPos = offset.top;
                currentNodes[i].positions[0] = xPos;
                currentNodes[i].positions[1] = yPos;
            },
            stop: function() {
                drawNodes();
            }
        });
        $(node).dblclick(function() {
            removeNode(currentNodes[i].value);
        });
        currentNodes[i].connections = currentNodes[i].connections || [];
        currentNodes[i].positions = currentNodes[i].positions || [0, 0];
        node.style.left = currentNodes[i].positions[0] + "px";
        node.style.top = currentNodes[i].positions[1] + "px";
        connections.push([node, currentNodes[i].connections]);
    }
    let ch = document.getElementById("draw").children;
    for (let i=0; i<connections.length; i++) {
        let start = $(connections[i][0]).position();
        let startWidth = $(connections[i][0]).width();
        let startHeight = $(connections[i][0]).height();
        console.log(connections[i]);
        for (let j=0; j<connections[i][1].length; j++) {
            let end = $("div[attr-nodeId=" + connections[i][1][j] + "]").position();
            let endWidth = $("div[attr-nodeId=" + connections[i][1][j] + "]").width();
            let endHeight = $("div[attr-nodeId=" + connections[i][1][j] + "]").height();
            $("#draw").line(start.left + (startWidth/2), start.top + (startHeight/2), end.left + (endWidth/2), end.top + (endHeight/2));
        }
    }
    $(".line").css("z-index", "90");

}

setInterval(updateProjects, 1000);