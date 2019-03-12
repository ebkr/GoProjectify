let loadedProject = null;
let currentNodes = [];
let oldNodes = [];
let drawnNodes = [];
let focusedNode = null;
let newNodePos = [0,0];
let isDragging = false;



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
    let builder = "";
    for (let i=0; i<splits.length; i++) {
        if (splits[i] == "<<GENERATE>>") {
            found = true;
        } else if (found) {
            builder += splits[i];
        }
    }
    return JSON.parse(builder);
}

function defineNodeInteraction(id, node, nodeArray) {
    $(node).draggable({
        scroll: false,
        containment: "#draw",
        drag: function () {
            isDragging = true;
            let thisPos = $(node).position();
            let x = thisPos.left;
            let y = thisPos.top;
            nodeArray.position[0] = x;
            nodeArray.position[1] = y;
        },
        stop: function () {
            isDragging = false;
            var oReq = new XMLHttpRequest();
            oReq.open("GET", "/api/LoadProject/Reposition/" + id + "/" + nodeArray.position[0] + ":" + nodeArray.position[1] + "/" + loadedProject);
            oReq.send();
            drawNodes();
            jsPlumb.repaintEverything();
        }
    });
    $(node).contextmenu(function (e) {
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
        $(link).click(function () {
            linkNodes(focusedNode);
        });

        let rename = document.createElement("button");
        $(rename).text("Rename");
        $(rename).click(function () {
            $(rename).text("To Implement");
        });

        let del = document.createElement("button");
        $(del).text("Delete");
        $(del).click(function () {
            removeNode(id);
        });

        $("#contextOptions").html(null);
        $("#contextOptions").append(rename);
        $("#contextOptions").append(link);
        $("#contextOptions").append(del);
        // Prevent browser context menu
        return false;
    });
}

function updateNode(id, nodeArray) {
    let node = null;
    drawnNodes.forEach(element => {
        if (element.getAttribute("attr-nodeId") == id) {
            node = element;
        }
    });
    if (node === null) {
        node = document.createElement("div");
        node.className = "node";
        node.setAttribute("attr-nodeId", id);
        node.setAttribute("draggable", "true");
        $("#draw").append(node);
    }
    node.innerText = nodeArray["value"].toString();
    node.style.left = nodeArray.position[0] + "px";
    node.style.top = nodeArray.position[1] + "px";
    return node;
}

// Draws generated nodes under {let currentNodes OF TYPE []}
function drawNodes() {
    let stringOfNew = currentNodes.toString();
    let stringOfOld = oldNodes.toString();
    // Check if nodes are different.
    if (stringOfOld !== stringOfNew) {

        idTracker = [];

        Object.keys(currentNodes).forEach(function(id) {
            idTracker.push(id);
            let result = updateNode(id, currentNodes[id]);
            drawnNodes = drawnNodes.filter(function(value, index, arr) {
                return value.getAttribute("attr-nodeId") !== id;
            });
            drawnNodes.push(result);
            defineNodeInteraction(id, result, currentNodes[id]);
        });

        drawnNodes.forEach(function(node, index, arr) {
            let found = false;
            idTracker.forEach(function(id) {
                if (node.getAttribute("attr-nodeId") == id) {
                    found = true;
                }
            });
            if (!found) {
                console.log("Missing node: " + node.innerHTML);
                arr.splice(index, 1);
                $(node).remove();
            }
        });

        jsPlumb.repaintEverything();
        
        jsPlumb.deleteEveryEndpoint();
    
        Object.keys(currentNodes).forEach(function(id) {
            jsPlumb.Defaults.Endpoints = ["Blank"];
            currentNodes[id].connections.forEach(con => {
                jsPlumb.connect({
                    source:$("div[attr-nodeId=" + id + "]"),
                    target:$("div[attr-nodeId=" + con + "]"),
                    anchors:["AutoDefault"],
                    connector:["Straight"],
                    overlays: ["PlainArrow"]
                });
            });    
        });
        
    }

    $(".jtk-endpoint").remove();
    $(".line").css("z-index", "90");
    // Do stuff
    currentNodes = oldNodes;
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

setInterval(function() {
    if (!isDragging && loadedProject !== null) {
        loadProject(loadedProject);
    }
}, 500);

setInterval(function() {
    if (isDragging) {
        jsPlumb.repaintEverything();
    }
}, 100);