using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using scripter.models.yamlFile;

namespace scripter_2.dtos
{
    public class YamlFileDto
    {
        public HeaderDto header { get; set; }

        public ConfigurationDto configuration { get; set; }

        public ActionDto action { get; set; }
        public ContextPlaceholderDto[] context_placeholders { get; set; }

        public object steps { get; set; }

        public object rules { get; set; }
    }
}
