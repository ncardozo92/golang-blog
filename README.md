# golang-blog

## Setup
### MySQL
docker run --name golang-blog-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql

### Configuración del firewall de windows para evitar bloqueo al correr la aplicación
El firewall de windows siempre preguntará por cada build que se corre si deseamos que ejecute 
ya que considera a la misma como una posible amenaza. Seguir estos pasos para impedir que 
siempre nos pregunte si queremos ejecutar la misma (Obtenido de un post de Stack Overflow):
<br>
1. Go to Windows Defender Firewall, in Left side menu you saw Inbound Rules click there, then Right Side menu you will see New Rule... click.
2. Choose Port opened from window -> Next Select TCP, then define which ports you want I choose 8080 click Next again, Choose Allow the connection Next, Check All Next, Give any Name Goland or anything you want and press Finish. Thats it