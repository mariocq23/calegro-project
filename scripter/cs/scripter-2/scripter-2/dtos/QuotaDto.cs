using Newtonsoft.Json;

namespace scripter_2.dtos
{
    public class QuotaDto
    {
        public string os { get; set; }

        public string cpu { get; set; }

        public string memory { get; set; }

        public string disk { get; set; }

        public string graphics { get; set; }
    }
}