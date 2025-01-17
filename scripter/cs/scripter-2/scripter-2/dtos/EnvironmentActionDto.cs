using Newtonsoft.Json;

namespace scripter_2.dtos
{
    public class EnvironmentActionDto
    {
        public DependencyDto dependencies { get; set; }

        public QuotaDto quota { get; set; }
    }
}