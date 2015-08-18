var Message = function(data) {
  data = data || {}
  this.user = m.prop(data.user)
  this.text = m.prop(data.text)
  this.ts = m.prop(data.ts)
};
Message.list = function() {
  var apiURL = "/log/" + m.route.param("channel") + "/" 
                       + m.route.param("year") + "/"
                       + m.route.param("month") + "/"
                       + m.route.param("day");
  return m.request({ method: "GET", url: apiURL }).then(function(result){
    return result.messages;
  });
}

var Log = {
  controller: function() {
    this.messages = Message.list();
    this.channel = m.route.param("channel");
    this.date = new Date(m.route.param("year"), m.route.param("month")-1, m.route.param("day"));
  },
  view: function(ctrl) {
    var navi = m("p", [
                 m.component(LogLink, { channel: ctrl.channel, date: ctrl.date, delta: -1 }),
                 m("a[href='/']", { config: m.route }, "チャンネル一覧へ戻る"),
                 m.component(LogLink, { channel: ctrl.channel, date: ctrl.date, delta: 1 })
        ]);
    return m("div", [
        navi,
        ctrl.messages().map(function(msg, index) {
          return m("div", { key: msg.ts }, [
                  m("img", { src: msg.user.icon }),
                  m("p", msg.user.name),
                  m("p", msg.ts),
                  m("p", msg.text)
          ]);
        }),
        navi
    ]);
  }
};

