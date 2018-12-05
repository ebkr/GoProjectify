## GoProjectify

-----------

### Goals

The software should produce a list of nodes that CAN be connected. It should be able to create a project, nodes within that project, links between nodes, and be able to perform the inverse of the previously mentioned requirements.

A node should not be able to be connect to another if it breaks the flow of the application, and recursive nodes should be impossible.

------------

### Achieved
-	#### Creatable
-	Projects
-	Nodes
-	Links

------------

### Remaining
-	Delete Links
-	Position Nodes

------------

#### Log
[2018/12/05]

Front-end has been implemented for the most part. I decided to not use any GUI framework, and instead rely on HTML. This proves beneficial as cross-platform is less of an issue. The only downside is that the project has to run as a web server.

The interactions are all the same, however the design of the projectify.go file had to be changed to incorporate the interface functionality.

Interactions are performed by sending a GET request to the project. (localhost:8080).

Nodes can be visually linked, and can display their links accordingly. Nodes can also be repositioned on the screen, although this doesn't yet save, positioning was an intended feature near the start of the project.

-----------

[2018/12/03]

Projects can be removed, and a list of all nodes can now be displayed.

The menu within loadCase() has been cleaned up, and it's now easier, as I don't need to add any new "case" values, or modify existing ones. I chose to use an array of functions to perform this, as I can continue to add and find cases easily, providing I keep them commented.

I've started the log pretty late, but most of the back-end is complete.

Once I've finished working on removing links, I'll need to rewrite the project slightly for linking and removing to work with node IDs.

Assuming it all works out fine, I'll need to implement a front-end to make it a usable system. I haven't decided which framework to use, however it's likely to be Electron due to the simplicity, however I may work on a Qt interface.
