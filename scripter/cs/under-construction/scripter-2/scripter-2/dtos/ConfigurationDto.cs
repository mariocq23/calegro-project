using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;
using scripter.models.yamlFile;

namespace scripter_2.dtos
{
    public class ConfigurationDto
    {
        public string user { get; set; }

        public string executor { get; set; }

        public string execution_mode { get; set; }

        public bool is_containerized { get; set; }

        public string public_password { get; set; }

        public string private_password_location { get; set; }

        public string certificate_location { get; set; }

        public string security { get; set; }

        public bool bypass_security { get; set; }

        public string location { get; set; }

        public ContextDto[] context { get; set; }

        public string encoding { get; set; }
    }
}
