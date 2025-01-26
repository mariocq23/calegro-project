using Newtonsoft.Json;

namespace scripter_2.dtos
{
    public class EnvironmentDto
    {
        public string name { get; set; }

        public EnvironmentActionDto[] action { get; set; }

        public string[] custom_properties { get; set; }
    }
}