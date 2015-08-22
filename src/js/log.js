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
    var navi = <p>
                <LogLink channel={ctrl.channel} date={ctrl.date} delta={-1} />
                <a href="/" config={m.route}>チャンネル一覧へ戻る</a>
                <LogLink channel={ctrl.channel} date={ctrl.date} delta={1} />
               </p>;
    return m("div", [
        navi,
        ctrl.messages().map(function(msg, index) {
          return <div className="log" key={msg.ts}>
                   <img src={msg.user.icon} />
                   <span className="name">{msg.user.name}</span>
                   <span className="time">{msg.ts}</span>
                   <span className="msg">{msg.text}</span>
                 </div>;
        }),
        navi
    ]);
  }
};

