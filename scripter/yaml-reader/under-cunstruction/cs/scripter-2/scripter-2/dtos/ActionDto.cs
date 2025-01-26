using Newtonsoft.Json;
using scripter.models.yamlFile;

namespace scripter_2.dtos
{
    public class ActionDto
    {
        public string name { get; set; }

        public string type { get; set; }

        public string main_function { get; set; }

        public string api { get; set; }

        public bool use_api { get; set; }

        public bool display_output_console { get; set; }

        public PlatformDto platform { get; set; }

        public string location { get; set; }

        public string action_executor { get; set; }

        public object initial_inputs { get; set; }

        public string[] environment_variables { get; set; }
    }
}