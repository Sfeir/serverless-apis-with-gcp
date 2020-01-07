import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.HttpServletRequest;
import java.io.PrintWriter;
import java.io.IOException;

public class Example {
    public void helloWorld(HttpServletRequest request, HttpServletResponse response) throws IOException {
        String param = request.getParameter("name");
        String name = param != null ? param : "World";
        response.getWriter().write("<h1>Cloud Functions</h1><h2>Hello " + name + "!</h2>");
    }
}
