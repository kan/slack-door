var LogLink = {
    view: function(ctrl, args) {
        var d = new Date(args.date);
        d.setDate(d.getDate() + args.delta);
        var apiURL = "/" + args.channel + "/"
                       + d.getFullYear() + "/"
                       + (d.getMonth() + 1) + "/"
                       + d.getDate();
        var label;
        if (args.label) {
            label = args.label;
        } else {
            label = d.getFullYear() + "/"
                       + (d.getMonth() + 1) + "/"
                       + d.getDate();
            if (args.delta > 0) {
                label = label + ">>";
            } else {
                label = "<<" + label;
            }
        }
        return m("a", {href: apiURL, config: m.route}, label);
    }
};
