var loadedProject = null;
var currentNodes = [];
var focusedNode = null;
var newNodePos = [0,0];


// Prompt new file creation
function createProject() {
    let projName = prompt("New Project Name");
    if (projName) {
        var oReq = new XMLHttpRequest();
        oReq.open("GET", "/api/NewProject/" + projName);
        oReq.send();
        loadedProject = projName;
    }
}

// Adds a new node
function addNode() {
    if (loadedProject) {
        let nodeName = prompt("New Node Name: ");
        if (nodeName) {
            var oReq = new XMLHttpRequest();
            oReq.open("GET", "/api/LoadProject/NewNode/" + nodeName + "/" + newNodePos[0] + ":" + newNodePos[1] + "/" + loadedProject);
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

// Link two nodes together
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

// Calls link(x,y). Creates link between two clicked nodes.
function linkNodes(obj) {
    $("#draw").attr("selected", "true");
    let clicked_a = obj;
    let clicked_b = false;
    $("#draw").click(function(e) {
        if (!clicked_a) {
            clicked_a = e.target;
        } else if (!clicked_b) {
            clicked_b = e.target;
            $("#draw").attr("selected", null);
            if (clicked_a != clicked_b) {
                if ($(clicked_a).attr("attr-nodeId") && $(clicked_b).attr("attr-nodeId")) {
                    link($(clicked_a).attr("attr-nodeId") , $(clicked_b).attr("attr-nodeId"));
                }
            }
        }
    })
}

// Removes a node from the project
function removeNode(nodeId) {
    if (loadedProject) {
        if (nodeId !== null) {
            var oReq = new XMLHttpRequest();
            oReq.open("GET", "/api/LoadProject/RemoveNode/" + nodeId + "/_/" + loadedProject);
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

// Load project details
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

// Update a list of files on the left
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

            if ($(node).text() == loadedProject) {
                $(node).attr("selected", true);
            }

            $(node).contextmenu(function() {
                $("#contextModal").attr("attr-active", true);
                $("#contextTitle").text("Manage Project: " + $(node).text())
                $("#contextOptions").html(null);

                let load = document.createElement("button");
                $(load).text("Load Project");
                $(load).click(function() {
                    loadProject($(node).text());
                });

                let del = document.createElement("button");
                $(del).text("Delete Project");
                $(del).click(function() {
                    var oReq = new XMLHttpRequest();
                    oReq.open("GET", "/api/LoadProject/DeleteProject/_/_/" + $(node).text());
                    oReq.send();
                    $(node).attr("selected", true);
                });

                $("#contextOptions").append(load);
                $("#contextOptions").append(del);

                return false;
            });
        }
    });
    oReq.open("GET", "/api/GetProjects");
    oReq.send();
}

// Turn <<GENERATE>> in to nodes.
function parseCallOutput(str, nodeTree) {
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
                node.connections.push(lineSplit[1]);
            } else if (lineSplit[0] == "Position") {
                node.positions = [Number(lineSplit[1]), Number(lineSplit[2])];
            }
        }
    }
    return nodeTree;
}

// Draws generated nodes under {let currentNodes OF TYPE []}
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
            containment: "#draw",
            drag: function(){
                var thisPos = $(node).position();
                var x = thisPos.left;
                var y = thisPos.top;
                currentNodes[i].positions[0] = x;
                currentNodes[i].positions[1] = y;
            },
            stop: function() {
                var oReq = new XMLHttpRequest();
                oReq.open("GET", "/api/LoadProject/Reposition/" + currentNodes[i].id + "/" + currentNodes[i].positions[0] + ":" + currentNodes[i].positions[1] + "/" + loadedProject);
                oReq.send();
                drawNodes();
            }
        });
        $(node).contextmenu(function(e) {
            e.stopPropagation();
            if (focusedNode) {
                $(focusedNode).attr("selected", null);
            }
            focusedNode = node;
            $(focusedNode).attr("selected", true);
            $("#contextModal").attr("attr-active", true);

            $("#contextTitle").text("Modifying Node: " + $(node).text());

            let link = document.createElement("button");
            $(link).text("Link");
            $(link).click(function() {
                linkNodes(focusedNode);
            });

            let rename = document.createElement("button");
            $(rename).text("Rename");
            $(rename).click(function() {
                $(rename).text("To Implement");
            });

            let del = document.createElement("button");
            $(del).text("Delete");
            $(del).click(function() {
                removeNode(currentNodes[i].id);
            });

            $("#contextOptions").html(null);
            $("#contextOptions").append(rename);
            $("#contextOptions").append(link);
            $("#contextOptions").append(del);
            // Prevent browser context menu
            return false;
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
        for (let j=0; j<connections[i][1].length; j++) {
            let end = $("div[attr-nodeId=" + connections[i][1][j] + "]").position();
            let endWidth = $("div[attr-nodeId=" + connections[i][1][j] + "]").width();
            let endHeight = $("div[attr-nodeId=" + connections[i][1][j] + "]").height();
            $("#draw").line(start.left + (startWidth/2), start.top + (startHeight/2), end.left + (endWidth/2), end.top + (endHeight/2));
        }
    }
    $(".line").css("z-index", "90");

}

// Constantly update list of projects
setInterval(updateProjects, 1000);

// Enable context escape
$("#contextModal").click(function() {
    $("#contextModal").attr("attr-active", null);
})

$("#draw").contextmenu(function(e) {
    e.stopPropagation();
    newNodePos[0] = e.pageX;
    newNodePos[1] = e.pageY;
    $("#contextModal").attr("attr-active", true);
    $("#contextTitle").text("What would you like to do?");
    $("#contextOptions").html(null);

    let newNode = document.createElement("button");
    $(newNode).text("Create a new node");
    $(newNode).click(function() {
        addNode();
    });
    $("#contextOptions").append(newNode);
    return false;
});