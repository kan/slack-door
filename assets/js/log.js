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
    var navi = {tag: "p", attrs: {}, children: [
                m.component(LogLink, {channel:ctrl.channel, date:ctrl.date, delta:-1}), 
                {tag: "a", attrs: {href:"/", config:m.route}, children: ["チャンネル一覧へ戻る"]}, 
                m.component(LogLink, {channel:ctrl.channel, date:ctrl.date, delta:1})
               ]};
    return m("div", [
        navi,
        ctrl.messages().map(function(msg, index) {
          return {tag: "div", attrs: {className:"log", key:msg.ts}, children: [
                   {tag: "img", attrs: {src:msg.user.icon}}, 
                   {tag: "span", attrs: {className:"name"}, children: [msg.user.name]}, 
                   {tag: "span", attrs: {className:"time"}, children: [msg.ts]}, 
                   {tag: "span", attrs: {className:"msg"}, children: [msg.text]}
                 ]};
        }),
        navi
    ]);
  }
};

