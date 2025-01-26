using Newtonsoft.Json;

namespace scripter_2.dtos
{
    public class DependencyDto
    {
        public string location { get; set; }

        public string[] list { get; set; }
    }
}