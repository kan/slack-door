var Channel = function(data) {
  data = data || {}
  this.id = m.prop(data.id)
  this.name = m.prop(data.name)
};
Channel.list = function() {
  return m.request({ method: "GET", url: "/channels" }).then(function(result){
    return result.channels;
  });
}

var ChannelList = {
  controller: function() {
    this.channels = Channel.list();
  },
  view: function(ctrl) {
    return m("ul", [
        ctrl.channels().map(function(channel, index) {
          return m("li", { key: channel.id },
                   m("a", {href: "/"+channel.name+"/2015/08/18", config: m.route}, channel.name));
        })
    ]);
  }
};

