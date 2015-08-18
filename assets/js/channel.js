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
                  m.component(LogLink, { channel: channel.name, date: new Date(), delta: 0, label: channel.name })
                  );
        })
    ]);
  }
};

