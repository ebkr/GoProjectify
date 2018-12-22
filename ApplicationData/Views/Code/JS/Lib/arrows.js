(function($) {

    let defaults = {
        color: "#000",
        thickness: "1px",
        strength: 0,
        parent: $(document.body)
    }

    $.fn.drawArrow = function(connectTo, options) {

        // Initialise Default Values
        let keys = Object.keys(defaults);
        if (options == null) {
            options = defaults;
        }
        for (let i=0; i<keys.length; i++) {
            if (options[keys[i]] == null || options[keys[i]] == undefined) {
                options[keys[i]] = defaults[keys[i]]
            }
        }

        // Draw Arrows
        let leftX = $(this).offset().left + ($(this).innerWidth() / 2);
        let leftY = $(this).offset().top + ($(this).innerHeight() / 2);

        let rightX = $(connectTo).offset().left + ($(connectTo).innerWidth() / 2);
        let rightY = $(connectTo).offset().top + ($(connectTo).innerHeight() / 2);

        let width = rightX - leftX;
        let height = rightY - leftY;

        let svg = document.createElementNS("http://www.w3.org/2000/svg", "path");

        let delta = width*options.strength;
        let leftHeightDelta = leftX + delta;
        let rightHeightDelta = rightX - delta;

        var path = "M"  + leftX + " " + leftY + " L" + rightX + " " + rightY;

        svg.setAttributeNS(null, "d", path);
        svg.setAttributeNS(null, "fill", "none");
        svg.setAttributeNS(null, "stroke", options.color);

        $(options.parent).append(svg)
    }

})($)